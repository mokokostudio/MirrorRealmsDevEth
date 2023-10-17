using TrueSync;

namespace MR.Battle {
    public class UnitSystem : BaseSystem<UnitCD> {
        public override string Group => "Update";
        public override int Order => 800;

        private static FP SPRate = Config.Battle.Constant.SPDeductRate / (FP)10000;

        protected override void Run() {
            Data.HPChange.Clear();
            if (Data.State == UnitState.Die)
                return;
            var anim = GetComponentData<UnitAnimCD>();
            if (Data.TryMove && Data.State != UnitState.Run && Data.State != UnitState.Dash && !Data.LimitMove)
                Data.State = UnitState.Run;
            if (!Data.TryMove && Data.State == UnitState.Run)
                Data.State = UnitState.Idle;
            if (!Data.TryMove && Data.State == UnitState.Dash)
                anim.CallAnim = "Sprint_Stop";
            if (Data.TryMove && Data.TryDash && anim.CurAnim == "Slide_F" && !Data.LimitMove) {
                anim.CallAnim = "";
                Data.State = UnitState.Dash;
            }
            if (!Data.AttackHolding && Data.TryAttack) {
                anim.CallCommond = PlayerCommond.Atk;
                anim.CallAtkRelease = false;
                Data.AttackHolding = true;
            }
            if (Data.AttackHolding && !Data.TryAttack) {
                anim.CallAtkRelease = true;
                Data.AttackHolding = false;
            }
            if (Data.TryDefense && !Data.LimitSkill) {
                anim.CallCommond = PlayerCommond.Defense;
                Data.TryDefense = false;
            }


            var list = Data.BattleGround.GetCollisionTarget(Entity);
            foreach (var e in list) {
                var box = e.GetComponentData<AttackBoxCD>();
                if (box != null && box.Unit != Data && box.Unit.Camp != Data.Camp) {
                    if (box.AtkLV == 1 && anim.Defense) {
                        box.Unit.GetComponentData<UnitAnimCD>().CallAnim = "Defense_Shock";
                        anim.CallAnim = "Defense_Success";
                        if (Data.Player != null)
                            Data.Player.SP = TSMath.Min(Data.Player.SP + Config.Battle.Constant.SPGrowbyDefense, Config.Battle.Constant.SPMax);
                    } else {
                        var boxAnim = box.Unit.GetComponentData<UnitAnimCD>();
                        var boxLoc = box.Unit.GetComponentData<LocationCD>();
                        anim.HitEffects.Add(new HitEffect {
                            attackPlayer = box.Unit.Player == null ? (byte)255 : box.Unit.Player.Index,
                            hitType = box.HitType,
                            hitLv = box.AtkLV,
                            hitDir = (boxAnim.LookTarget ? boxAnim.Face : boxLoc.Face) + (box.Rotation.eulerAngles.y + 180) * TSMath.Deg2Rad,
                            hitSP = box.SP * SPRate,
                            hitHP = Data.BattleGround.Rage ? (box.HP * 2) : box.HP
                        });
                        anim.Pause = box.Paush;
                        boxAnim.Pause = box.Paush;
                    }
                }
            }
        }
    }
}
