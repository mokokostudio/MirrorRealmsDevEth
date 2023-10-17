using System.Collections.Generic;
using TrueSync;
using UnityEngine;

namespace MR.Battle {
    [AttachTo(typeof(UnitCD))]
    public class UnitDisplayCD : ComponentData {
        public Vector3 Offset { get; set; }
        public Vector3 PositionDelta { get; set; }
        public float FaceDeg { get; set; }
        public float RealFaceDeg { get; set; }
        public GameObject Go { get; set; }
        public UnitDisplay UnitDisplay { get; set; }
        public string Posture { get; set; }
        public UnitState LastState { get; set; }
        public string LastAnimMode { get; set; }
        public bool NewAnim { get; set; }
        public List<GameObject> WeaponGos { get; } = new List<GameObject>();
        public string RunMode { get; set; }
        public int DefLv { get; set; }
        public bool OnHit { get; set; }
        public FP Speed { get; set; }
    }
}
