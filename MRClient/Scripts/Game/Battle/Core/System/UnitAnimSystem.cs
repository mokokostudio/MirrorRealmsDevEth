using System.Collections.Generic;
using TrueSync;
using UnityEngine.SocialPlatforms.Impl;
using UnityEngine.UIElements;

namespace MR.Battle {
    public class UnitAnimSystem : BaseSystem<UnitAnimCD> {
        public override string Group => "Update";
        public override int Order => 600;
        public FP Interval => Battle.Instance.Interval;
        private UnitCD m_Unit;
        private LocationCD m_Location;
        private BattleGroundScoreCD m_Score;

        protected override void Run() {
            Data.AnimOffset = TSVector.zero;

            m_Unit = GetComponentData<UnitCD>();
            m_Location = GetComponentData<LocationCD>();
            m_Score = m_Unit.BattleGround.GetComponentData<BattleGroundScoreCD>();

            if (Data.Pause > 0) {
                Data.Pause -= 1;
                Data.AnimSpeed = .01f;
            } else
                Data.AnimSpeed = 1f;

            Data.AnimTime += Interval * Data.AnimSpeed;

            if (Data.AnimClip != null) {
                foreach (var cs in Data.AnimClip.comboSkillPoints) {
                    if (Across(cs.startPoint))
                        if (cs.release) {
                            if (Data.CallAtkRelease)
                                Data.CallAnim = cs.skill;
                            else
                                Data.AtkReleaseNext = cs.skill;
                        } else
                            Data.AtkCastNext = cs.skill;
                    if (Across(cs.endPoint))
                        if (cs.release) {
                            if (cs.endPoint == Data.AnimClip.length && Data.CallCommond == PlayerCommond.None)
                                Data.CallAnim = cs.skill;
                            Data.AtkReleaseNext = null;
                        } else
                            Data.AtkCastNext = null;
                }
                if (!Data.AnimClip.moveable && m_Unit.State == UnitState.Run)
                    StopAnim();
            }

            if (m_Unit.State != UnitState.Die) {
                if (Data.CallUnarm && !m_Unit.LimitDodge) {
                    TurnCommon();
                    Data.CallUnarm = false;
                }

                switch (Data.CallCommond) {
                    case PlayerCommond.Atk:
                        if (Data.AtkCastNext != null && !Data.CallUnarm) {
                            Data.CallAnim = Data.AtkCastNext;
                            Data.CallCommond = PlayerCommond.None;
                            Data.CallAtkRelease = false;
                        } else if (!m_Unit.LimitSkill) {
                            if (Data.Mode != "Common") {
                                switch (m_Unit.State) {
                                    case UnitState.Run:
                                        if (BattleResources.GetAnimationDataClip(Data.Posture, Data.Mode, "Atk_R") != null)
                                            Data.CallAnim = "Atk_R";
                                        else
                                            Data.CallAnim = "Atk";
                                        break;
                                    case UnitState.Dash:
                                        if (BattleResources.GetAnimationDataClip(Data.Posture, Data.Mode, "Atk_D") != null)
                                            Data.CallAnim = "Atk_D";
                                        else if (BattleResources.GetAnimationDataClip(Data.Posture, Data.Mode, "Atk_R") != null) {
                                            Data.CallAnim = "Atk_R";
                                            m_Unit.State = UnitState.Run;
                                        } else
                                            Data.CallAnim = "Atk";
                                        break;
                                    default:
                                        Data.CallAnim = "Atk";
                                        break;
                                }
                            }
                            Data.CallCommond = PlayerCommond.None;
                        }
                        break;
                    case PlayerCommond.Defense:
                        if (!m_Unit.LimitSkill) {
                            if (Data.Mode != "Common" && Data.Mode != "Bow")
                                Data.CallAnim = "Defense_Cast";
                            Data.CallCommond = PlayerCommond.None;
                            Data.CallAtkRelease = false;
                        }
                        break;
                    case PlayerCommond.Dodge:
                        if (!m_Unit.LimitDodge) {
                            m_Unit.Player.EP -= Config.Battle.Constant.EpReductByDodge;
                            Data.CallAnim = "Slide_F";
                            Data.CallCommond = PlayerCommond.None;
                            Data.CallAtkRelease = false;
                        }
                        break;
                }

                if (Data.CallAtkRelease) {
                    if (Data.AtkReleaseNext != null) {
                        Data.CallAnim = Data.AtkReleaseNext;
                        Data.CallCommond = PlayerCommond.None;
                    }
                }

                Data.OnHit = false;
                while (Data.HitEffects.Count > 0) {
                    Data.OnHit = true;
                    var he = Data.HitEffects[0];

                    if (Data.DefLv <= he.hitLv) {
                        switch (he.hitType) {
                            case HitType.None:
                                ReduceHP(he);
                                if (m_Unit.HP <= 0) {
                                    TurnCommon();
                                    Data.CallAnim = "Die_FL";
                                    Die(he);
                                }
                                break;
                            case HitType.Normal:
                                m_Location.Face = he.hitDir;
                                if (he.hitLv == 0 && Data.Defense || Data.Mode == "Common" || Data.Mode == "Bow" || he.hitLv > Data.DefLv || Data.CurAnim != null && Data.CurAnim.StartsWith("Hit")) {
                                    ReduceHP(he);
                                    if (m_Unit.HP <= 0) {
                                        TurnCommon();
                                        Data.CallAnim = "Die_FL";
                                        Die(he);
                                    } else
                                        Data.CallAnim = "Hit_F";
                                } else {
                                    Data.CallAnim = "Defense_Hit";
                                    if (m_Unit.Player != null)
                                        m_Unit.Player.SP -= he.hitSP;
                                }
                                break;
                            case HitType.Large:
                                m_Location.Face = he.hitDir;
                                TurnCommon();
                                if (m_Unit.Player != null)
                                    m_Unit.Player.SP = 0;
                                ReduceHP(he);
                                if (m_Unit.HP <= 0) {
                                    Data.CallAnim = "Die_Large_B";
                                    Die(he);
                                } else
                                    Data.CallAnim = "Hit_Large";
                                break;
                        }
                    }
                    Data.HitEffects.RemoveAt(0);
                }
            }

            if (Data.CallAnim == "") {
                if (Data.AnimClip != null)
                    StopAnim();
            } else if (Data.CallAnim != null) {
                if (Data.CallAnim == "Slide_F")
                    m_Location.Face = m_Unit.TryFace;
                if (Data.CallAnim == "Sprint_Stop" && GetAnim("Sprint_End_L") != null && GetAnim("Sprint_End_R") != null) {
                    var sprint = GetAnim("Sprint");
                    var length = sprint.length;
                    var animTime = Data.AnimTime % length;
                    if (animTime < length / 4 || animTime > length * 3 / 4)
                        Data.CallAnim = "Sprint_End_R";
                    else
                        Data.CallAnim = "Sprint_End_L";
                }
                StopAnim();
                PlayAnim(Data.CallAnim);
                if (Data.AnimClip == null) {
                    UnityEngine.Debug.LogError($"Anim Error, Mode{Data.Mode},Anim:{Data.CallAnim}");
                } else {
                    if (Data.AnimClip.skillPoint > 0)
                        m_Unit.LimitSkill = true;
                    if (Data.AnimClip.movePoint > 0)
                        m_Unit.LimitMove = true;
                    if (Data.AnimClip.dodgePoint > 0)
                        m_Unit.LimitDodge = true;
                    Data.CurAnim = Data.CallAnim;
                }
                Data.CallAnim = null;
            }

            if (Data.AnimClip != null) {
                if (Across(Data.AnimClip.weaponPoint)) {
                    var equip = GetComponentData<UnitEquipCD>();
                    Data.Mode = Config.Equips.WeaponType[equip.Weapon.Type].Mode;
                    equip.Complete = true;
                }
                if (Across(Data.AnimClip.skillPoint))
                    m_Unit.LimitSkill = false;
                if (Across(Data.AnimClip.movePoint))
                    m_Unit.LimitMove = false;
                if (Across(Data.AnimClip.dodgePoint))
                    m_Unit.LimitDodge = false;
                for (int i = 0; i < Data.AnimClip.attackBoxConfigs.Length; i++) {
                    var box = Data.AnimClip.attackBoxConfigs[i];
                    if (Across(box.startPoint)) {
                        var entity = World.CreateEntity();
                        var boxcd = entity.AddComponentData<AttackBoxCD>();
                        boxcd.BattleGround = m_Unit.BattleGround;
                        boxcd.Unit = m_Unit;
                        boxcd.Position = box.position;
                        boxcd.Rotation = TSMatrix.CreateRotationY(box.rotation.y) * TSMatrix.CreateRotationX(box.rotation.x) * TSMatrix.CreateRotationZ(box.rotation.z);
                        boxcd.Scale = box.scale;
                        boxcd.AtkLV = box.atkLv;
                        boxcd.HitType = box.hitType;
                        boxcd.SP = box.sp;
                        boxcd.HP = box.hp;
                        boxcd.Paush = box.paush;
                        Data.AttackBoxEntities.Add(i, entity);
                    }
                    if (Across(box.endPoint)) {
                        var e = Data.AttackBoxEntities[i];
                        m_Unit.BattleGround.UnregisterPhysic(e);
                        Data.AttackBoxEntities.Remove(i);
                        World.RemoveEntity(e);
                    }
                }

                for (int i = 0; i < Data.AnimClip.bulletConfigs.Length; i++) {
                    var conf = Data.AnimClip.bulletConfigs[i];
                    if (Across(conf.startPoint)) {
                        //var entity = World.CreateEntity();
                        //entity = entity.AddComponentData<BulletCD>();
                    }
                }

                for (int i = 0; i < Data.AnimClip.endureConfigs.Length; i++) {
                    var endure = Data.AnimClip.endureConfigs[i];
                    if (Across(endure.startPoint))
                        Data.DefLv = endure.defLV;
                    if (Across(endure.endPoint))
                        Data.DefLv = 0;
                }

                for (int i = 0; i < Data.AnimClip.faceConfigs.Length; i++) {
                    var face = Data.AnimClip.faceConfigs[i];
                    if (Across(face.point)) {
                        if (face.lockTarget)
                            UpdateTarget();
                        if (face.maxAngle > 0)
                            m_Location.Face = GetFace(face.maxAngle);
                    }
                }

                for (int i = 0; i < Data.AnimClip.defenseConfigs.Length; i++) {
                    var defense = Data.AnimClip.defenseConfigs[i];
                    if (Across(defense.startPoint))
                        Data.Defense = true;
                    if (Across(defense.endPoint))
                        Data.Defense = false;
                }

                for (int i = 0; i < Data.AnimClip.runModeConfigs.Length; i++) {
                    var conf = Data.AnimClip.runModeConfigs[i];
                    if (Across(conf.startPoint))
                        Data.RunMode = conf.mode;
                    if (Across(conf.endPoint))
                        Data.RunMode = null;
                }

                for (int i = 0; i < Data.AnimClip.lookTargetConfigs.Length; i++) {
                    var conf = Data.AnimClip.lookTargetConfigs[i];
                    if (Across(conf.startPoint))
                        Data.LookTarget = true;
                    if (Across(conf.endPoint))
                        Data.LookTarget = false;
                }



                for (int i = 0; i < Data.AnimClip.spConfigs.Length; i++) {
                    var conf = Data.AnimClip.spConfigs[i];
                    if (Across(conf.startPoint))
                        if (m_Unit.Player != null)
                            m_Unit.Player.SP -= conf.sp;
                }

                Data.Impulse = default;
                for (int i = 0; i < Data.AnimClip.impulseConfigs.Length; i++) {
                    var conf = Data.AnimClip.impulseConfigs[i];
                    if (Across(conf.startPoint))
                        if (!conf.onlySelf || m_Unit.Player != null && m_Unit.Player == Battle.Instance.CameraPlayer)
                            Data.Impulse = conf.power;
                }

                for (int i = 0; i < Data.AnimClip.effectConfigs.Length; i++) {
                    var conf = Data.AnimClip.effectConfigs[i];
                    if (Across(conf.startPoint))
                        Data.Fxs.Add(CreateEffectFromConfig(conf));
                    if (Across(conf.endPoint))
                        Data.Fxs.Remove(CreateEffectFromConfig(conf));
                }

                Data.AnimOffset = GetAnimMove(Data.AnimClip.moveData, Data.AnimTime * 60, Data.AnimTime * 60 + 2) * Data.AnimSpeed;

                if (Data.LookTarget) {
                    Data.Face = GetFace(360, false);
                    if (m_Unit.State == UnitState.Idle)
                        m_Location.Face = Data.Face;
                }

                if (Data.AnimTime >= Data.AnimClip.length) {
                    if (Data.CurAnim.StartsWith("Die"))
                        AddComponentData<EntityRemoveCD>();
                    else {
                        if (!Data.AnimClip.moveable)
                            m_Unit.State = UnitState.Idle;
                        StopAnim();
                    }
                }
            }
        }

        private FX CreateEffectFromConfig(CharacterAnimationDataClip.EffectConfig conf) {
            var effect = new FX();
            effect.start = conf.startPoint;
            effect.end = conf.endPoint;
            effect.target = conf.target;
            effect.follow = conf.follow;
            effect.speed = conf.speed;
            effect.position = conf.position;
            effect.rotation = conf.rotation;
            effect.scale = conf.scale;
            return effect;
        }

        private void ReduceHP(HitEffect he) {
            m_Unit.HP -= he.hitHP;
            m_Unit.HPChange.Add(-he.hitHP);
            m_Score.AddOutput(he.attackPlayer, he.hitHP);
            if (m_Unit.Player != null)
                m_Score.AddTake(m_Unit.Player.Index, he.hitHP);
        }

        private void Die(HitEffect he) {
            if (m_Unit.Player != null)
                m_Score.Kill(he.attackPlayer, m_Unit.Player.Index);
            var equip = m_Unit.GetComponentData<UnitEquipCD>();
            if (equip != null)
                RemoveComponentData(equip);
        }

        private void TurnCommon() {
            if (Data.Mode != "Common") {
                Data.Mode = "Common";
                StopAnim();
            }
        }

        private FP GetFace(FP limit, bool withRun = true) {
            FP target;
            if (m_Unit.Target != null) {
                var delta = m_Unit.Target.Position - m_Location.Position;
                target = TSMath.Atan2(delta.x, delta.z);
            } else if (withRun && m_Unit.TryMove)
                target = m_Unit.TryFace;
            else if (m_Unit.Player != null)
                target = m_Unit.Player.CamDir;
            else
                return default;
            var deltaAngle = TSMath.DeltaAngle(m_Location.Face * TSMath.Rad2Deg, target * TSMath.Rad2Deg);
            if (deltaAngle > 0 && deltaAngle > limit)
                deltaAngle = limit;
            if (deltaAngle < 0 && deltaAngle < -limit)
                deltaAngle = -limit;
            return TSMath.DeltaAngle(0, m_Location.Face * TSMath.Rad2Deg + deltaAngle) * TSMath.Deg2Rad;
        }

        private void UpdateTarget() {
            var targets = new List<LocationCD>();
            foreach (var e in World.GetEntities<UnitCD>()) {
                var unit = e.GetComponentData<UnitCD>();
                if (unit == m_Unit || unit.Camp == m_Unit.Camp)
                    continue;
                targets.Add(e.GetComponentData<LocationCD>());
            }
            LocationCD target = null;
            if (m_Unit.TryMove)
                target = GetTarget(targets, m_Unit.TryFace);
            else if (m_Unit.Target != null)
                return;
            else if (m_Unit.Player != null)
                target = GetTarget(targets, m_Unit.Player.CamDir);
            if (target != null)
                m_Unit.Target = target;
        }

        private FP s_MaxLockAngle = Config.Battle.Constant.MaxLockAngle;
        private FP s_MaxLockDistance = Config.Battle.Constant.MaxLockDistance / (FP)100;

        private LocationCD GetTarget(List<LocationCD> locations, FP face) {
            LocationCD target = null;
            var weight = FP.Zero;
            foreach (var location in locations) {
                var dv = location.Position - m_Location.Position;
                var dis = dv.magnitude;
                if (dis > s_MaxLockDistance)
                    continue;
                var angle = TSMath.Abs(TSMath.DeltaAngle(face * TSMath.Rad2Deg, TSMath.Atan2(dv.x, dv.z) * TSMath.Rad2Deg));
                if (angle > s_MaxLockAngle)
                    continue;
                var w = 2 - dis / s_MaxLockDistance - angle / s_MaxLockAngle;
                if (w > weight) {
                    weight = w;
                    target = location;
                }
            }
            return target;
        }

        private void StopAnim() {
            Data.AnimTime = 0;
            Data.AnimClip = null;
            Data.CurAnim = null;
            m_Unit.LimitSkill = false;
            m_Unit.LimitMove = false;
            m_Unit.LimitDodge = false;
            Data.AtkCastNext = null;
            Data.AtkReleaseNext = null;
            Data.Defense = false;
            Data.RunMode = null;
            Data.LookTarget = false;
            Data.Fxs.Clear();
            foreach (var e in Data.AttackBoxEntities.Values) {
                m_Unit.BattleGround.UnregisterPhysic(e);
                World.RemoveEntity(e);
            }
            Data.AttackBoxEntities.Clear();
            Data.DefLv = 0;
        }

        private TSVector GetAnimMove(TSVector2[] data, FP from, FP to) {
            var start = TSMath.Clamp(TSMath.Floor(from), 0, data.Length).AsInt();
            var end = TSMath.Clamp(TSMath.Ceiling(to), 0, data.Length).AsInt();
            var sum = TSVector2.zero;
            for (int i = start; i < end; i++) {
                var step = data[i];
                FP space = FP.One;
                if (i < from)
                    space -= from - i;
                if (i > end)
                    space -= i - end;
                sum += step * space;
            }
            return new TSVector(sum.x, 0, sum.y);
        }

        private bool Across(FP c) {
            var a = Data.AnimTime;
            var b = Data.AnimTime + Interval * Data.AnimSpeed;
            return a <= c && b > c;
        }

        private void PlayAnim(string name) {
            Data.AnimClip = GetAnim(name);
            if (name.StartsWith("Die"))
                m_Unit.State = UnitState.Die;
            else if (Data.AnimClip && !Data.AnimClip.moveable)
                m_Unit.State = UnitState.Skill;
        }

        private CharacterAnimationDataClip GetAnim(string name) {
            return BattleResources.GetAnimationDataClip(Data.Posture, Data.Mode, name);
        }
    }
}
