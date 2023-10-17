using System.IO;

namespace MR.Battle {
    public class OperateDodgeEnd : IOperateData {
        public bool Equals(IOperateData other) => other is OperateDodgeEnd;

        public void Read(BinaryReader reader) { }
        public void Write(BinaryWriter writer) { }
    }
}
