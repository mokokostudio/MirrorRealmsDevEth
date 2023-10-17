using System;

namespace MR.Battle {
    public interface ISystem {
        World World { get; set; }
        Type Type { get; }
        string Group { get; }
        int Order { get; }
        void Run(object data);
    }
}
