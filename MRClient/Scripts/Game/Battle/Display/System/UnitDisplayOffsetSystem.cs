using BDFramework.UFlux;
using Unity.VisualScripting;
using UnityEngine;

namespace MR.Battle {
    public class UnitDisplayOffsetSystem : BaseSystem<UnitDisplayCD> {
        public override string Group => "Update";
        public override int Order => 101;

        protected override void Run() {
            var offset = GetComponentData<MoveCD>().Offset;
            Data.Offset = offset.ToVector();
            if (Data.Go != null) {
                var disP = Data.Go.transform.position;
                var location = Data.GetComponentData<LocationCD>();
                var orgP = location.Position.ToVector();
                Data.PositionDelta = orgP - disP;
            }
        }
    }
    public class UnitDisplayDamageTextSystem : BaseSystem<UnitDisplayCD> {
        public override string Group => "Update";
        public override int Order => 101;
        protected override void Run() {
            var unit = GetComponentData<UnitCD>();
            foreach (var damage in unit.HPChange) {
                var damageGo = Object.Instantiate(UFluxUtils.LoadAsync<GameObject>("Logic/Damage.prefab"), Data.Go.transform, false);
                damageGo.transform.localPosition = new Vector3(0, 1.6f, 0);
                damageGo.GetComponent<DamageText>().Damage = damage.AsInt();
            }
        }
    }
}
