using System.IO;

namespace MR.Battle {
    public class OperateMoveStop : IOperateData {
        public bool Equals(IOperateData other) => other is OperateMoveStop;
        public void Read(BinaryReader reader) { }
        public void Write(BinaryWriter writer) { }
    }
}
