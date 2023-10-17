using System.IO;
using TrueSync;

namespace MR.Battle {
    public class OperateCameraMove : IOperateData {
        public FP toward;
        public bool Equals(IOperateData other) => other is OperateCameraMove && (other as OperateCameraMove).toward == toward;
        public void Read(BinaryReader reader) => toward = FP.FromRaw(reader.ReadInt64());
        public void Write(BinaryWriter writer) => writer.Write(toward.RawValue);
    }
}
