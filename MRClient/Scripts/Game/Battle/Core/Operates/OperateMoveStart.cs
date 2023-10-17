using System.IO;

namespace MR.Battle {
    public class OperateMoveStart : IOperateData {
        public bool Equals(IOperateData other) => other is OperateMoveStart;

        public void Read(BinaryReader reader) { }
        public void Write(BinaryWriter writer) { }
    }
}
