using System.IO;
using TrueSync;

namespace MR.Battle {
    public class OperateMoveUpdate : IOperateData {
        public FP toward;
        public bool Equals(IOperateData other) => other is OperateMoveUpdate && (other as OperateMoveUpdate).toward == toward;
        public void Read(BinaryReader reader) => toward = FP.FromRaw(reader.ReadInt64());
        public void Write(BinaryWriter writer) => writer.Write(toward.RawValue);
    }
}
