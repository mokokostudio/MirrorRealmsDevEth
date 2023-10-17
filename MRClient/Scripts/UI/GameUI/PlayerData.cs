using MR.Net.Proto.Battle;

public static class PlayerData {
    public static string Name { get; private set; }
    public static WeaponPB[] Weapons { get; set; }
    public static int[] Equiped { get; set; }

    public static void Set(UserInfo data) {
        Name = data.Name;
        Weapons = data.Weapons;
        Equiped = data.Equiped;
    }
}
