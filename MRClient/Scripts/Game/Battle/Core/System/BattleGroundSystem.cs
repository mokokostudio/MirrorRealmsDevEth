using TrueSync.Physics3D;
using TSWorld = TrueSync.Physics3D.World;

namespace MR.Battle {
    public class BattleGroundSystem : BaseSystem<BattleGroundCD> {
        public override string Group => "Update";
        public override int Order => 1000;
        protected override void Run() {
            if (Data.OperateData != null) {
                var fData = Data.OperateData;
                Data.OperateData = null;
                int offset = 4;
                while (offset < fData.Length) {
                    var player = fData[offset++];
                    var list = Data.OperateDataProcesser.Deserialize(fData, ref offset);
                    foreach (var op in list)
                        Data.Players[player].OperateDatas.Enqueue(op);
                }
            }
            if (Data.WatchNext) {
                Data.WatchNext = false;
                var Idx = Data.PlayerWatch.Index;
                do {
                    Idx = (byte)((Idx + 1) % 6);
                } while (!Data.Players.ContainsKey(Idx));
                Data.PlayerWatch = Data.Players[Idx];
            }
            if (Data.TSWorld == null) {
                Data.TSWorld = new TSWorld(new CollisionSystemBrute());
                Data.TSWorld.Events.BodiesBeginCollide += c => {
                    Data.collisionMap.Add(c.body1, c.Body2);
                    Data.collisionMap.Add(c.Body2, c.Body1);
                };
            }
            Data.collisionMap.Clear();
            Data.TSWorld.Step(Battle.Instance.Interval);
        }
    }
}
