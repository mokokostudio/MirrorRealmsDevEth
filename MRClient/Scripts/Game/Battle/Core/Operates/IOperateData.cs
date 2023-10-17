using System;
using System.IO;

namespace MR.Battle {
    public interface IOperateData : IEquatable<IOperateData> {
        void Read(BinaryReader reader);
        void Write(BinaryWriter writer);
    }
}
