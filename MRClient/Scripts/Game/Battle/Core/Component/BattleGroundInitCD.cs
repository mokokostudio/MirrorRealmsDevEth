using System.Collections.Generic;
using TrueSync;

namespace MR.Battle {
    public class BattleGroundInitCD : ComponentData {
        public int[] PlayerIndex { get; set; }
        public int HandleIndex { get; set; }
        public int Sead { get; set; }
        public TSVector[] BornPoints { get; set; }
        public Dictionary<int, ulong[]> WeaponDic { get; set; }
    }
}
