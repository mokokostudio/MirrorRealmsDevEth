using System;
using System.Collections.Generic;
using TrueSync;

namespace MR.Battle {
    public class PlayerCD : ComponentData {
        public BattleGroundCD BattleGround { get; set; }
        public byte Index { get; set; }
        public Random Random { get; set; }
        public TSVector BornPos { get; set; }
        public FP BornFace { get; set; }
        public Camp Camp { get; set; }

        public Config.Equips.WeaponData[] Weapons { get; set; }
        public Queue<IOperateData> OperateDatas { get; } = new Queue<IOperateData>();
        public UnitCD Unit { get; set; }
        public FP CamDir { get; set; }
        public int SelectCard { get; set; }

        public FP SP { get; set; }
        public int SpGrow { get; set; }

        public FP EP { get; set; }

        public PlayerState State { get; set; }
        public FP StateTime { get; set; }
        public List<Card> CardPool { get; } = new List<Card>();
        public CardResult[] HandCard { get; } = new CardResult[] { new CardResult(), new CardResult(), new CardResult() };

        public class Card {
            public Config.Equips.WeaponCardType CardType { get; set; }
            public int WeaponIdx { get; set; }
        }

        public class CardResult {
            public List<Card> Cards { get; } = new List<Card>();
        }
    }
}
