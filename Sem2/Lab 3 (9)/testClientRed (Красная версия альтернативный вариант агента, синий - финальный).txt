import socket
import json
import math
import time
import random
import os

HOST = "127.0.0.1"
PORT = 8080
TEAM = "red"
STATE_INGAME = "InGame"

LOG_BUFFER = []
LAST_FLUSH_TIME = time.time()

def light_log(msg):
    global LAST_FLUSH_TIME
    line = f"[{time.strftime('%H:%M:%S')}] {msg}"
    LOG_BUFFER.append(line)
    if len(LOG_BUFFER) > 60 or (time.time() - LAST_FLUSH_TIME > 10):
        flush_logs()

def flush_logs():
    global LAST_FLUSH_TIME
    if not LOG_BUFFER: return
    try:
        with open(f"log_{TEAM}.txt", "a", encoding="utf-8") as f:
            f.write("\n".join(LOG_BUFFER) + "\n")
        LOG_BUFFER.clear()
        LAST_FLUSH_TIME = time.time()
    except Exception: pass

BASE_DROP_RADIUS    = 1.9    
GOLEM_CHASE_SPEED   = 3.5    
GOLEM_VISION_DIST   = 10.0   
GOLEM_LOSE_DIST     = 15.0   
GOLEM_FOV_COS       = 0.766  
PICKUP_RADIUS       = 1.3    
DROP_WEIGHT_THRESH  = 3.6    
ENEMY_STEAL_DROP    = 1.8    
ENEMY_STEAL_PANIC   = 6.0    
STEAL_CHASE_DIST    = 6.0    
STEAL_BONUS_TIME    = 5.0    

BASE_OFFSETS = [
    {"x":  1.2, "z":  1.2}, {"x": -1.2, "z":  1.2},
    {"x":  1.2, "z": -1.2}, {"x": -1.2, "z": -1.2},
    {"x":  0,   "z":  0}
]

team_memory = {
    "base_pos":       None,
    "assignments":    {},    
    "golem_prev_pos": {},
    "golem_dir":      {},    
    "prev_points":    0,
    "enemy_points":   0,
    "stats": {
        "steal_attempts": 0, "steal_success": 0, "stuns": 0, "stucks": 0,
        "last_summary": 0.0, "tactics": {}, "total_ticks": 0
    },
    "agent_prev_stun": {}
}

EXPLORE_POINTS = [{"x":80,"z":80}, {"x":-80,"z":80}, {"x":80,"z":-80}, {"x":-80,"z":-80}, {"x":0,"z":0}]

class Agent:
    def __init__(self, agent_id, idx=0):
        self.id             = str(agent_id)
        self.idx            = idx
        self.target_id      = None
        self.last_cmd_time  = 0.0
        self.last_sent_pos      = None
        self.last_sent_pos_time = 0.0
        self.pos_hist       = []
        self.unstuck_target = None
        self.unstuck_until  = 0.0
        self.explore_index  = idx
        self.steal_bonus_until = 0.0
        self.last_steal_pct = 0.0
        self.current_tactic = "IDLE"
        self.last_action    = None
        self.pickup_attempts = 0
        self.blacklisted_treasures = {} 

    @staticmethod
    def dist(a, b):
        if not a or not b: return 1e9
        return math.hypot(a["x"] - b["x"], a["z"] - b["z"])

    def _cmd(self, action, target=None, oid=None):
        now = time.time()
        if action in ("pickup", "drop", "steal"): self.pos_hist.clear()
        
        is_priority = (self.current_tactic == "DELIVERY" or "STEAL" in self.current_tactic)
        throttle = 0.2 if is_priority else 0.4

        if action == "position" and target is not None:
            if (self.last_sent_pos and self.dist(self.last_sent_pos, target) < 0.3 and now - self.last_sent_pos_time < throttle):
                return None
            self.last_sent_pos, self.last_sent_pos_time = target, now

        if action not in ("pickup", "drop", "steal", "ready") and now - self.last_cmd_time < 0.08:
            return None
        self.last_cmd_time = now
        
        if action in ("pickup", "drop", "steal"):
            if self.last_action != action:
                light_log(f"Agent {self.id} -> {action.upper()} {f'({oid})' if oid else ''}")
                self.last_action = action
        else: self.last_action = None
        
        cmd = {"Id": self.id, "action": action}
        if target is not None: cmd["target"] = target
        return cmd

    def _release_assignment(self):
        if self.target_id and team_memory["assignments"].get(self.target_id) == self.id:
            del team_memory["assignments"][self.target_id]
        self.target_id = None

    def decide(self, me, treasures, bases, game_state, golems, all_agents):
        pos = me.get("pos")
        if not pos or game_state != STATE_INGAME: return None
        
        has_treasure, my_weight = me.get("hasTreasure", False), me.get("weight", 0)
        base_pos = team_memory["base_pos"] or {"x": 0, "y": 0, "z": 0}
        now = time.time(); d_base = self.dist(pos, base_pos)
        
        cur_pct = me.get("stealChargePercentage", 0) or 0
        if (self.last_steal_pct >= 80 and cur_pct < 30 and has_treasure):
            self.steal_bonus_until = now + STEAL_BONUS_TIME
            team_memory["stats"]["steal_success"] += 1
            light_log(f"!!! [STEAL SUCCESS] Agent {self.id} !!!")
        self.last_steal_pct = cur_pct

        if me.get("isStunned"):
            self.current_tactic = "STUNNED"; self._release_assignment()
            self.pos_hist.clear(); self.steal_bonus_until = 0.0
            return None

        off = BASE_OFFSETS[self.idx % len(BASE_OFFSETS)]
        my_base_target = {"x": base_pos["x"] + off["x"], "y": base_pos["y"], "z": base_pos["z"] + off["z"]}

        if now < self.unstuck_until and self.unstuck_target:
            self.current_tactic = "UNSTUCK"
            return self._cmd("position", self.unstuck_target)

        self.pos_hist.append({"x": pos["x"], "z": pos["z"]})
        stuck_thresh = 50 if d_base < 3.0 else 25
        if len(self.pos_hist) > stuck_thresh:
            self.pos_hist.pop(0)
            xs, zs = [p["x"] for p in self.pos_hist], [p["z"] for p in self.pos_hist]
            if (max(xs)-min(xs) < 0.6) and (max(zs)-min(zs) < 0.6):
                if d_base < 3.0 and has_treasure: pass
                else:
                    team_memory["stats"]["stucks"] += 1; self.unstuck_until = now + 1.2
                    light_log(f"[STUCK FIX] Agent {self.id} at d={d_base:.1f}m, dash AWAY")
                    if d_base < 5.0:
                        vx, vz = pos["x"]-base_pos["x"], pos["z"]-base_pos["z"]
                        ln = math.hypot(vx, vz) or 1.0
                        self.unstuck_target = {"x": pos["x"]+vx/ln*10, "y":0, "z": pos["z"]+vz/ln*10}
                    else:
                        ang = random.uniform(0, 2*math.pi)
                        self.unstuck_target = {"x": pos["x"]+math.cos(ang)*10, "y":0, "z": pos["z"]+math.sin(ang)*10}
                    self.pos_hist.clear()
                    return self._cmd("position", self.unstuck_target)
        elif len(self.pos_hist) > 60: self.pos_hist.pop(0)

        base_eff = 10.0 if (team_memory["prev_points"] - team_memory["enemy_points"]) >= -200 else 7.0
        eff_dist = max(base_eff, 8.5 if my_weight > 100 else 0)
        
        danger_golem, threat_level, chasing_us = self._assess_golem_threat(pos, golems, eff_dist)

        if danger_golem and threat_level >= 2:
            self.current_tactic = "FLEE_GOLEM"
            gpos = danger_golem["pos"]
            gd = self.dist(pos, gpos)
            if threat_level == 3 and chasing_us:
                if has_treasure:
                    if now > self.steal_bonus_until and max(5.0*(1.0-0.005*my_weight), 0.5) < DROP_WEIGHT_THRESH:
                        self._release_assignment(); return self._cmd("drop")
                    return self._cmd("position", my_base_target)
                else:
                    vx, vz = pos["x"]-gpos["x"], pos["z"]-gpos["z"]
                    ln = math.hypot(vx, vz) or 1.0
                    return self._cmd("position", {"x": pos["x"]+vx/ln*18, "y":0, "z": pos["z"]+vz/ln*18})

        if has_treasure:
            self.current_tactic = "DELIVERY"; self._release_assignment()
            if d_base < BASE_DROP_RADIUS: return self._cmd("drop")
            return self._cmd("position", my_base_target)

        if me.get("stealAbilityReady"):
            best_e, best_sc = None, -1.0
            for e in all_agents:
                if str(e.get("team","")).lower() == TEAM or not e.get("hasTreasure") or e.get("isStunned"): continue
                sc = (1.0 + e.get("weight",0)*0.015) / (self.dist(pos, e["pos"]) + 1.0)
                if sc > best_sc: best_sc, best_e = sc, e
            if best_e:
                ed = self.dist(pos, best_e["pos"])
                if ed < 1.0: 
                    team_memory["stats"]["steal_attempts"] += 1; return self._cmd("steal")
                if ed < 6.0: 
                    self.current_tactic = "STEAL_CHASE"; return self._cmd("position", best_e["pos"])

        self.blacklisted_treasures = {k: v for k, v in self.blacklisted_treasures.items() if v > now}
        available = [t for t in treasures if not t.get("isPicked") and t.get("holderAgentId") is None and t["id"] not in self.blacklisted_treasures]
        if not available:
            self.current_tactic = "EXPLORE"
            ep = EXPLORE_POINTS[self.idx % len(EXPLORE_POINTS)]
            if self.dist(pos, ep) < 3.0: self.explore_index += 1
            return self._cmd("position", ep)

        def score(t):
            d_to, d_h = self.dist(pos, t["pos"]), self.dist(t["pos"], base_pos)
            s = (t.get("value", 0) or 10) / max((d_to + d_h) / (max(1.0-0.005*t.get("weight", 0), 0.1)), 0.1)
            if team_memory["assignments"].get(t["id"]) not in (self.id, None): s *= 0.002
            if d_to < 2.0: s *= 80.0
            return s

        best = max(available, key=score)
        if self.target_id and self.target_id != best["id"]:
            curr = next((t for t in available if t["id"] == self.target_id), None)
            if curr and self.dist(pos, curr["pos"]) < 15.0 and score(best) < score(curr) * 1.5: best = curr
        
        self.target_id = best["id"]; team_memory["assignments"][self.target_id] = self.id
        self.current_tactic = "GATHERING"
        if self.dist(pos, best["pos"]) < PICKUP_RADIUS:
            self.pickup_attempts += 1
            if self.pickup_attempts > 15:
                self.blacklisted_treasures[best["id"]] = now + 20.0
                light_log(f"!! [PICKUP FAIL] Blacklisting {best['id']}")
                self._release_assignment(); self.pickup_attempts = 0
                return None
            return self._cmd("pickup", best["pos"], oid=best["id"])
        self.pickup_attempts = 0
        return self._cmd("position", best["pos"])

    def _assess_golem_threat(self, pos, golems, eff):
        best_g, best_t, our_f = None, 0, False
        for g in golems:
            gd = self.dist(pos, g["pos"])
            gs, gt = g.get("state", "Patrol"), g.get("targetAgentId")
            t, f = 0, False
            if gs == "Chase":
                if gt == self.id and gd < GOLEM_LOSE_DIST: t, f = 3, True
                elif gd < 5.0: t = 2
            else:
                g_dir = team_memory["golem_dir"].get(g["id"], {"x": 0, "z": 0})
                dl = math.hypot(g_dir["x"], g_dir["z"])
                if dl > 0.01 and gd < eff:
                    to_x, to_z = pos["x"]-g["pos"]["x"], pos["z"]-g["pos"]["z"]
                    tl = math.hypot(to_x, to_z) or 1.0
                    if (g_dir["x"]/dl*to_x/tl + g_dir["z"]/dl*to_z/tl) > GOLEM_FOV_COS: t, f = 1, True
            if t > best_t: best_t, best_g, our_f = t, g, f
        return best_g, best_t, our_f

def print_summary():
    s, r, b = team_memory["stats"], team_memory.get('prev_points',0), team_memory.get('enemy_points',0)
    res = "WIN" if r > b else ("LOSS" if r < b else "DRAW")
    light_log(f"--- ИТОГ: RED {r} vs BLUE {b} [{res}] | Станы: {s['stuns']}, Застрял: {s['stucks']}, Кражи: {s['steal_success']} ---")
    if s["total_ticks"] > 0:
        parts = [f"{k}: {v/s['total_ticks']*100:.1f}%" for k, v in s["tactics"].items() if v > 0]
        light_log("Тактики: " + " | ".join(parts))

def start_client():
    sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    sock.settimeout(None); agents_ai, agent_idx = {}, 0
    try:
        sock.connect((HOST, PORT))
        light_log(f"[SYSTEM] Подключились: {TEAM.upper()}")
        sock.sendall((json.dumps({"team": TEAM}) + "\n").encode("utf-8"))
        buffer, player_id = "", None
        while True:
            chunk = sock.recv(65536)
            if not chunk: break
            buffer += chunk.decode("utf-8", errors="ignore")
            while "\n" in buffer:
                line, buffer = buffer.split("\n", 1); msg = json.loads(line)
                mtype = msg.get("type")
                if mtype == "joinAccepted":
                    player_id = msg.get("playerId")
                    sock.sendall((json.dumps({"actions":[{"action":"ready"}]}) + "\n").encode("utf-8"))
                elif mtype == "gameEvent":
                    ev = msg.get("eventType")
                    if ev == "start":
                        team_memory["assignments"].clear(); team_memory["prev_points"] = 0; agents_ai.clear(); agent_idx = 0
                        team_memory["stats"].update({"steal_attempts":0, "steal_success":0, "stuns":0, "stucks":0, "tactics":{}, "total_ticks":0, "last_summary": time.time()})
                    elif ev == "result": print_summary(); flush_logs()
                if not player_id or "agents" not in msg: continue
                if team_memory["base_pos"] is None:
                    for b in msg["bases"]:
                        if b["team"].lower() == TEAM: team_memory["base_pos"] = b["pos"]
                for b in msg["bases"]:
                    if b["team"].lower() == TEAM: team_memory["prev_points"] = b["points"]
                    else: team_memory["enemy_points"] = b["points"]
                for a in msg["agents"]:
                    if str(a.get("team","")).lower() == TEAM:
                        aid = a["agentId"]
                        if a.get("isStunned") and not team_memory["agent_prev_stun"].get(aid): team_memory["stats"]["stuns"] += 1
                        team_memory["agent_prev_stun"][aid] = a.get("isStunned")
                for g in msg.get("golems", []):
                    gid, prev = g["id"], team_memory["golem_prev_pos"].get(g["id"])
                    if prev:
                        dx, dz = g["pos"]["x"]-prev["x"], g["pos"]["z"]-prev["z"]
                        ln = math.hypot(dx, dz)
                        if ln > 0.01:
                            nx, nz = dx/ln, dz/ln
                            old = team_memory["golem_dir"].get(gid, {"x": nx, "z": nz})
                            team_memory["golem_dir"][gid] = {"x": old["x"]*0.6+nx*0.4, "z": old["z"]*0.6+nz*0.4}
                    team_memory["golem_prev_pos"][gid] = g["pos"]
                my_agents = [a for a in msg["agents"] if str(a.get("team","")).lower() == TEAM]
                my_agents.sort(key=lambda a: 0 if a.get("hasTreasure") else 1)
                actions = []
                for a in my_agents:
                    aid = a["agentId"]
                    if aid not in agents_ai: agents_ai[aid] = Agent(aid, agent_idx); agent_idx += 1
                    cmd = agents_ai[aid].decide(a, msg.get("treasures", []), msg.get("bases", []), "InGame", msg.get("golems", []), msg["agents"])
                    tac = agents_ai[aid].current_tactic
                    team_memory["stats"]["tactics"][tac] = team_memory["stats"]["tactics"].get(tac, 0) + 1
                    team_memory["stats"]["total_ticks"] += 1
                    if cmd: actions.append(cmd)
                if actions: sock.sendall((json.dumps({"actions":actions}) + "\n").encode("utf-8"))
                if time.time() - team_memory["stats"]["last_summary"] > 30:
                    team_memory["stats"]["last_summary"] = time.time(); print_summary(); flush_logs()
    except Exception as e: light_log(f"ERROR: {e}")
    finally: sock.close(); flush_logs()

if __name__ == "__main__": start_client()
