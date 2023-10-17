using System;
using TrueSync;

namespace MR.Battle {
    public class BattleGroundInitSystem : BaseSystem<BattleGroundInitCD> {
        public override string Group => "Update";
        public override int Order => int.MaxValue;
        protected override void Run() {
            var bg = Data.Entity.GetComponentData<BattleGroundCD>();
            bg.Random = new Random(Data.Sead);
            var processer = bg.OperateDataProcesser = new OperateDataProcesser();
            processer.RegistTwin<OperateMoveStart, OperateMoveStop>();
            processer.Regist<OperateMoveUpdate>();
            processer.Regist<OperateCameraMove>();
            processer.Regist<OperateSelectCard>();
            processer.Regist<OperateDiscardWeapon>();
            processer.RegistTwin<OperateDodgeStart, OperateDodgeEnd>();
            processer.RegistTwin<OperateAttackStart, OperateAttackEnd>();
            processer.Regist<OperateDefense>();
            var score = GetComponentData<BattleGroundScoreCD>();
            for (int i = 0; i < Data.PlayerIndex.Length; i++) {
                var idx = (byte)Data.PlayerIndex[i];
                var entity = World.CreateEntity();
                var player = entity.AddComponentData<PlayerCD>();
                score.AddPlayer(idx);
                player.BattleGround = bg;
                player.Index = idx;
                player.Random = new Random(bg.Random.Next());
                player.Camp = idx < 3 ? Camp.PlayerA : Camp.PlayerB;
                var a = idx * 60 * TSMath.Deg2Rad;
                player.BornPos = Data.BornPoints[idx];
                player.BornFace = TSMath.Atan2(-player.BornPos.z, player.BornPos.x) - TSMath.Pi / 2;
                player.Weapons = new Config.Equips.WeaponData[3];
                for (int j = 0; j < 3; j++)
                    player.Weapons[j] = Config.Equips.Weapon[Data.WeaponDic[idx][j]];
                bg.Players.Add(idx, player);
                if (Data.HandleIndex == -1 && i == 0)
                    bg.PlayerWatch = player;
                if (Data.HandleIndex == idx)
                    bg.PlayerHandle = player;
            }
            RemoveComponentData(Data);
        }
    }
}
