using System.Collections;
using TrueSync;

namespace MR.Battle {
    public class PlayerSystem : BaseSystem<PlayerCD> {
        public override string Group => "Update";
        public override int Order => 900;
        private static FP s_Interval = Battle.Instance.Interval;
        private static FP s_EPStep = Battle.Instance.Interval * Config.Battle.Constant.EPMax / Config.Battle.Constant.EpRecoverTime;

        protected override void Run() {
            if (Data.Unit == null || Data.Unit.Entity == null)
                Born();
            if (Data.Unit.State != UnitState.Die) {
                Operate();
                StateLogic();
                Data.EP = TSMath.Min(Data.EP + s_EPStep, Config.Battle.Constant.EPMax);
            }
        }

        private void Born() {
            var entity = World.CreateEntity();
            Data.Unit = entity.AddComponentData<UnitCD>();
            Data.Unit.Player = Data;
            Data.Unit.BattleGround = Data.BattleGround;
            Data.Unit.MoveSpeed = 4;
            Data.Unit.HP = Config.Battle.Constant.HPMax;
            Data.Unit.Camp = Data.Camp;
            var location = Data.Unit.GetComponentData<LocationCD>();
            location.Position = Data.BornPos;
            location.Face = Data.BornFace;
            Data.CamDir = Data.BornFace;
            if (Data.BattleGround.PlayerHandle == Data)
                Data.BattleGround.PreCameraDir = Data.BornFace;
            Data.EP = Config.Battle.Constant.EPMax;
            Data.State = PlayerState.WaitCard;
            Data.StateTime = 0;

            //entity = World.CreateEntity();
            //var Unit = entity.AddComponentData<UnitCD>();
            //Unit.BattleGround = Data.BattleGround;
            //Unit.MoveSpeed = 4;
            //Unit.HP = Config.Battle.Constant.HPMax;
        }

        private void Operate() {
            while (Data.OperateDatas.TryDequeue(out var i)) {
                switch (i) {
                    case OperateMoveStart:
                        Data.Unit.TryMove = true;
                        break;
                    case OperateMoveStop:
                        Data.Unit.TryMove = false;
                        break;
                    case OperateMoveUpdate op:
                        Data.Unit.TryFace = op.toward;
                        break;
                    case OperateCameraMove op:
                        Data.CamDir = op.toward;
                        break;
                    case OperateSelectCard op:
                        Data.SelectCard = op.index;
                        break;
                    case OperateDiscardWeapon:
                        Data.SP = 0;
                        break;
                    case OperateDodgeStart:
                        if (Data.EP >= Config.Battle.Constant.EpReductByDodge) {
                            Data.Unit.GetComponentData<UnitAnimCD>().CallCommond = PlayerCommond.Dodge;
                            Data.Unit.TryDash = true;
                        }
                        break;
                    case OperateDodgeEnd:
                        Data.Unit.TryDash = false;
                        break;
                    case OperateAttackStart:
                        if (Data.State == PlayerState.Fight)
                            Data.Unit.TryAttack = true;
                        break;
                    case OperateAttackEnd:
                        if (Data.State == PlayerState.Fight)
                            Data.Unit.TryAttack = false;
                        break;
                    case OperateDefense:
                        if (Data.State == PlayerState.Fight)
                            Data.Unit.TryDefense = true;
                        break;
                }
            }
        }

        private void StateLogic() {
            switch (Data.State) {
                case PlayerState.WaitCard:
                    Data.StateTime -= s_Interval;
                    if (Data.StateTime <= 0) {
                        if (Data.CardPool.Count == 0)
                            for (int i = 0; i < 3; i++)
                                for (int j = 0; j < 5; j++)
                                    Data.CardPool.Add(new PlayerCD.Card { CardType = Data.Weapons[i].Cards[j], WeaponIdx = i });
                        for (int i = 0; i < 5; i++) {
                            int idx = Data.Random.Next(Data.CardPool.Count);
                            var card = Data.CardPool[idx];
                            Data.HandCard[card.WeaponIdx].Cards.Add(card);
                            Data.CardPool.Remove(card);
                        }
                        Data.State = PlayerState.ShowHand;
                        Data.SelectCard = -1;
                    }
                    break;
                case PlayerState.ShowHand:
                    if (Data.SelectCard != -1) {
                        Data.State = PlayerState.WaitFightStart;
                        Data.Unit.AddComponentData<UnitEquipCD>().Weapon = Data.Weapons[Data.SelectCard];
                        Data.SpGrow = Config.Battle.Constant.SPGrowBase + Data.HandCard[Data.SelectCard].Cards.Count * Config.Battle.Constant.SPGrowFromWeapon;
                    }
                    break;
                case PlayerState.WaitFightStart:
                    if (Data.Unit.GetComponentData<UnitEquipCD>() == null) {
                        Data.State = PlayerState.Fight;
                        Data.SP = Data.SpGrow;
                    }
                    break;
                case PlayerState.Fight:
                    if (Data.SP <= 0) {
                        Data.State = PlayerState.WaitFightEnd;
                        Data.Unit.GetComponentData<UnitAnimCD>().CallUnarm = true;
                        for (int i = 0; i < 3; i++)
                            Data.HandCard[i].Cards.Clear();
                    }
                    break;
                case PlayerState.WaitFightEnd:
                    if (Data.Unit.GetComponentData<UnitAnimCD>().Mode == "Common") {
                        Data.Unit.TryAttack = false;
                        Data.Unit.TryDefense = false;
                        Data.State = PlayerState.WaitCard;
                        Data.StateTime = Config.Battle.Constant.DualDuration;
                    } 
                    break;
            }
        }
    }
}
