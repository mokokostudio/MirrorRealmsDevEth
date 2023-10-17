using BDFramework.UFlux;
using TrueSync;
using UnityEngine;
using UnityEngine.Animations;

namespace MR.Battle {
    public class UnitDisplaySystem : BaseSystem<UnitDisplayCD> {
        public override string Group => "Display";
        public override int Order => 800;
        public float Interval => Battle.Instance.Interval.AsFloat();
        private UnitCD m_Unit;
        private static float s_Side = 24.5f;

        protected override void Run() {
            var remove = GetComponentData<EntityRemoveCD>();
            if (remove != null) {
                if (Data.Go != null) {
                    Object.Destroy(Data.Go);
                    Data.Go = null;
                }
                return;
            }
            m_Unit = GetComponentData<UnitCD>();
            var location = GetComponentData<LocationCD>();
            var anim = GetComponentData<UnitAnimCD>();
            var equip = GetComponentData<UnitEquipCD>();
            var unitDis = Data;
            if (Data.Go == null) {
                Data.Go = m_Unit.Player.Index < 3 ? BattleResources.GetAvatar("Eos") : BattleResources.GetAvatar("EosBlack");
                Data.Go.transform.position = location.Position.ToVector();
                Data.RealFaceDeg = Data.FaceDeg = location.Face.AsFloat() * Mathf.Rad2Deg;
                Data.Go.transform.eulerAngles = new Vector3(0, Data.FaceDeg, 0);
                Data.UnitDisplay = Data.Go.AddComponent<UnitDisplay>();
                Data.Posture = "Girl";
                Data.UnitDisplay.Init(Data.Posture, () => GetMoveSpeed(m_Unit), () => GetAnimFace(unitDis));

                var hpbarGo = Object.Instantiate(UFluxUtils.LoadAsync<GameObject>("Logic/HPBar.prefab"), Data.Go.transform, false);
                hpbarGo.transform.localPosition = new Vector3(0, 1.5f, 0);
                hpbarGo.GetComponent<HPBar>().Unit = m_Unit;
            }
            var pos = Data.Go.transform.position;
            var face = (anim.LookTarget ? anim.Face : location.Face).AsFloat() * Mathf.Rad2Deg;

            Data.RealFaceDeg += Mathf.DeltaAngle(Data.RealFaceDeg, location.Face.AsFloat() * Mathf.Rad2Deg) * Time.deltaTime * 10;

            pos += Data.Offset / Interval * Time.deltaTime;

            var ds = Vector3.MoveTowards(Vector3.zero, Data.PositionDelta, Time.deltaTime);
            Data.PositionDelta -= ds;
            var ds2 = Data.PositionDelta * Time.deltaTime * 10;
            Data.PositionDelta -= ds2;
            pos += ds + ds2;

            pos.x = Mathf.Clamp(pos.x, -s_Side, s_Side);
            pos.z = Mathf.Clamp(pos.z, -s_Side, s_Side);

            Data.Go.transform.position = pos;

            if ((m_Unit.State == UnitState.Run || m_Unit.State == UnitState.Dash) && !anim.LookTarget)
                Data.FaceDeg += Mathf.DeltaAngle(Data.FaceDeg, face) * Time.deltaTime * 2;
            else
                Data.FaceDeg = Mathf.MoveTowardsAngle(Data.FaceDeg, face, Time.deltaTime * 720);
            Data.Go.transform.eulerAngles = new Vector3(0, Data.FaceDeg, 0);

            if (Data.LastAnimMode != anim.Mode) {
                Data.UnitDisplay.SetMode(anim.Mode);

                Data.LastAnimMode = anim.Mode;
            }

            if (anim.AnimTime == 0 && anim.AnimClip != null) {
                if (Data.NewAnim) {
                    if (m_Unit.State != UnitState.Idle && m_Unit.State != UnitState.Run && m_Unit.State != UnitState.Dash)
                        Data.UnitDisplay.Play(anim.AnimClip, false);
                    else
                        Data.UnitDisplay.Play(anim.AnimClip, true);
                    Data.NewAnim = false;
                }
            } else
                Data.NewAnim = true;

            if (Data.LastState != m_Unit.State) {
                switch (m_Unit.State) {
                    case UnitState.Hit:
                    case UnitState.Skill:
                        Data.UnitDisplay.PreIdle();
                        break;
                    case UnitState.Idle:
                        Data.UnitDisplay.Idle();
                        break;
                    case UnitState.Run:
                        Data.UnitDisplay.Move();
                        break;
                    case UnitState.Dash:
                        Data.UnitDisplay.Dash();
                        break;
                }
                Data.LastState = m_Unit.State;
            }

            if (equip != null && equip.Complete && Data.WeaponGos.Count == 0) {
                var weaponType = Config.Equips.WeaponType[equip.Weapon.Type];
                for (int i = 0; i < equip.Weapon.Prefabs.Count; i++) {
                    var wGo = BattleResources.GetWeapon(equip.Weapon.Prefabs[i]);
                    var wp = Data.Go.transform.Find($"Weapons/{i + 1}");
                    var pc = wp.GetComponent<ParentConstraint>();
                    pc.weight = 0;
                    wGo.transform.SetParent(wp, false);
                    for (int j = 0; j < pc.sourceCount; j++) {
                        var source = pc.GetSource(j);
                        source.weight = weaponType.EquipPosition[i] == j ? 1 : 0;
                        pc.SetSource(j, source);
                    }
                    pc.weight = 1;
                    Data.WeaponGos.Add(wGo);
                }
            }

            if (anim.Mode == "Common" && Data.WeaponGos.Count > 0) {
                foreach (var go in Data.WeaponGos)
                    FreeItem.Free(go, 5);
                Data.WeaponGos.Clear();
            }

            if (anim.RunMode != Data.RunMode) {
                Data.RunMode = anim.RunMode;
                Data.UnitDisplay.SetMove(anim.RunMode);
            }

            if (anim.DefLv != Data.DefLv) {
                Data.DefLv = anim.DefLv;
                Data.UnitDisplay.SetDefRim(anim.DefLv);
            }

            if (anim.OnHit != Data.OnHit) {
                Data.OnHit = anim.OnHit;
                Data.UnitDisplay.SetSplash();
            }

            if (anim.Impulse != default) {
                Data.UnitDisplay.Impulse(anim.Impulse);
                anim.Impulse = default;
            }

            if (anim.AnimSpeed != Data.Speed) {
                Data.Speed = anim.AnimSpeed;
                Data.UnitDisplay.SetSpeed(Data.Speed.AsFloat());
            }

            Data.UnitDisplay.UpdateEffect(anim.Fxs);
        }

        private float GetMoveSpeed(UnitCD unit) {
            var speed = unit.MoveSpeed.AsFloat();
            if (m_Unit.State == UnitState.Dash)
                speed *= 1.5f;
            return speed;
        }

        private float GetAnimFace(UnitDisplayCD unitDis) {
            var deg = Mathf.Repeat(unitDis.FaceDeg - unitDis.RealFaceDeg, 360f) / 360f;
            return deg;
        }
    }
}
