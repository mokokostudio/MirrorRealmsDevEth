namespace MR.Battle {
    public class UnitEquipSystem : BaseSystem<UnitEquipCD> {
        public override string Group => "Update";
        public override int Order => 790;

        protected override void Run() {
            var unit = GetComponentData<UnitCD>();
            var anim = GetComponentData<UnitAnimCD>();
            if (!unit.LimitSkill)
                anim.CallAnim = $"Equip_{Config.Equips.WeaponType[Data.Weapon.Type].Mode}";
            if (Data.Complete)
                RemoveComponentData(Data);
        }
    }
}
