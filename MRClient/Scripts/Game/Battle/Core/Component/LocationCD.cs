using TrueSync;

namespace MR.Battle {
    [AttachTo(typeof(UnitCD))]
    public class LocationCD : ComponentData {
        public TSVector Position { get; set; }
        public FP Face { get; set; }
    }
}
