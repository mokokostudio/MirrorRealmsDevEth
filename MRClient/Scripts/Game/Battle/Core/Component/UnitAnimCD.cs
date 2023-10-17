using System;
using System.Collections.Generic;
using TrueSync;
using UnityEngine;

namespace MR.Battle {
    [AttachTo(typeof(UnitCD))]
    public class UnitAnimCD : ComponentData {
        public CharacterAnimationDataClip AnimClip { get; set; }
        public string Mode { get; set; } = "Common";
        public string Posture { get; set; } = "Girl";
        public string RunMode { get; set; }
        public UnitState LastState { get; set; }
        public FP AnimTime { get; set; }
        public TSVector AnimOffset { get; set; }
        public string CallAnim { get; set; }
        public string CurAnim { get; set; }
        public PlayerCommond CallCommond { get; set; }
        public bool CallAtkRelease { get; set; }
        public bool CallUnarm { get; set; }
        public string AtkCastNext { get; set; }
        public string AtkReleaseNext { get; set; }
        public Dictionary<int, Entity> AttackBoxEntities { get; } = new Dictionary<int, Entity>();
        public bool LookTarget { get; set; }
        public FP Face { get; set; }

        public List<HitEffect> HitEffects { get; } = new List<HitEffect>();

        public int DefLv { get; set; }
        public bool Defense { get; set; }

        public bool OnHit { get; set; }

        public TSVector Impulse { get; set; }

        public int Pause { get; set; }
        public FP AnimSpeed { get; set; } = 1;

        public List<FX> Fxs { get; } = new List<FX>();
    }

    public struct FX : IEquatable<FX> {
        public FP start;
        public FP end;
        public GameObject target;
        public bool follow;
        public float speed;
        public Vector3 position;
        public Vector3 rotation;
        public Vector3 scale;

        public bool Equals(FX other) {
            return start == other.start
                && end == other.end
                && follow == other.follow
                && speed == other.speed
                && position == other.position
                && rotation == other.rotation
                && scale == other.scale;
        }
    }

    public struct HitEffect {
        public byte attackPlayer;
        public HitType hitType;
        public int hitLv;
        public FP hitDir;
        public FP hitSP;
        public FP hitHP;
    }

    public enum PlayerCommond {
        None = 0,
        Atk,
        Defense,
        Dodge
    }
}
