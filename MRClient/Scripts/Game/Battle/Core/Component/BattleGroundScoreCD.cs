using System;
using System.Collections.Generic;
using TrueSync;

namespace MR.Battle {
    [AttachTo(typeof(BattleGroundCD))]
    public class BattleGroundScoreCD : ComponentData {
        public Dictionary<byte, ScoreData> Datas { get; } = new Dictionary<byte, ScoreData>();
        public List<KillData> KillRecord { get; } = new List<KillData>();
        public int killA;
        public int killB;

        public void AddPlayer(byte index) {
            if (Datas.ContainsKey(index))
                return;
            Datas[index] = new ScoreData();
        }

        public void Kill(byte killer, byte killed) {
            if (Datas.TryGetValue(killer, out var data)) {
                data.kill++;
                if (killer < 3)
                    killA++;
                else
                    killB++;
                Datas[killer] = data;
            }
            if (Datas.TryGetValue(killed, out var data2)) {
                data2.die++;
                Datas[killed] = data2;
            }
            KillRecord.Add(new KillData { killer = killer, killed = killed });
        }

        public void AddOutput(byte index, FP damage) {
            if (Datas.TryGetValue(index, out var data)) {
                data.output+= damage;
                Datas[index] = data;
            }
        }

        public void AddTake(byte index, FP damage) {
            if (Datas.TryGetValue(index, out var data)) {
                data.take+= damage;
                Datas[index] = data;
            }
        }

        public struct ScoreData {
            public int kill;
            public int die;
            public FP output;
            public FP take;
        }

        public struct KillData {
            public int killer;
            public int killed;
        }
    }
}
