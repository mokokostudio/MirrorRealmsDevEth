namespace MR.Battle {
    public class UnitEquipCD : ComponentData {
        public Config.Equips.WeaponData Weapon { get; set; }
        public bool Complete { get; set; }
    }
}
