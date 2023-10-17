using TrueSync;

namespace MR.Battle {
    [AttachTo(typeof(UnitCD))]
    public class MoveCD : ComponentData {
        public TSVector Offset { get; set; }
    }
}
