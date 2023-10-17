using System.IO;

namespace MR.Battle {
    public class OperateDodgeStart : IOperateData {
        public bool Equals(IOperateData other) => other is OperateDodgeStart;

        public void Read(BinaryReader reader) { }
        public void Write(BinaryWriter writer) { }
    }
}
