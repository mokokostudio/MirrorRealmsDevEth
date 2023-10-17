using BDFramework.UFlux;
using Cinemachine;
using MR.Battle;
using System;
using System.Collections;
using System.Collections.Generic;
using System.ComponentModel.Composition.Hosting;
using TrueSync;
using UnityEngine;
using UnityEngine.Playables;

public class UnitDisplay : MonoBehaviour {
    private enum State {
        Idle,
        Move,
        Dash
    }
    private AuAnimator2 m_Animator2;
    private string m_Posture;
    private string m_Mode;

    private AnimationClip m_IdleClip;
    private RunData m_RunGroup;
    private RunData m_DashGroup;

    private State m_State;

    private List<Material> m_Materials = new List<Material>();

    private Color m_DefRimLasColor;
    private Color m_DefRimCurColor;
    private Color m_DefRimTarColor;
    private float m_DefRimColorLerp;

    private float m_SplashLerp;
    private float m_SplashValue;

    private float m_Speed;

    private Dictionary<FX, PlayableDirector> m_FXDirectors = new Dictionary<FX, PlayableDirector>();
    private List<FX> m_FXRecord = new List<FX>();
    private GameObject m_DashEffect;
    private Func<float> m_FaceGet;

    public void Init(string posture, Func<float> moveSpeedGet, Func<float> faceGet) {
        m_Posture = posture;
        m_FaceGet = faceGet;
        SetMode("Common");
        m_Animator2 = new AuAnimator2(name, GetComponent<Animator>(), m_IdleClip, m_RunGroup.clips, m_RunGroup.moveDistance, moveSpeedGet, faceGet);
        SetMatirals();
        m_DashEffect = Instantiate(UFluxUtils.LoadAsync<GameObject>("Effects/FX_Common_Sprint_Loop.prefab"), transform, false);
        m_DashEffect.SetActive(false);
    }

    private void Update() {
        var step = Time.deltaTime * m_Speed;
        m_Animator2.Evaluate(step);
        foreach (var kv in m_FXDirectors) {
            var fx = kv.Key;
            var director = kv.Value;
            director.time += step * fx.speed;
            director.Evaluate();
        }

        UpdateRim();
        UpdateSplash();
        if (m_DashEffect.activeInHierarchy)
            m_DashEffect.transform.localEulerAngles = new Vector3(0, -m_FaceGet() * 360, 0);
    }

    private void OnAnimatorMove() { }

    private void OnDestroy() {
        m_Animator2.Dispose();
    }

    public void Play(CharacterAnimationDataClip data, bool blendMove) {
        //Debug.Log($"PlayAnim:{data.name},blendmove:{blendMove}");
        if (data.clip != null)
            m_Animator2.PlayAnim(data.clip, blendMove);
        else
            m_Animator2.PlayAnim(data.customAnims, blendMove);
        m_DashEffect.SetActive(false);
    }

    public void PreIdle() {
        m_State = State.Idle;
    }

    public void Idle() {
        m_State = State.Idle;
        m_Animator2.PlayIdle();
    }

    public void Move() {
        m_State = State.Move;
        m_Animator2.SetMove(m_RunGroup.clips, m_RunGroup.moveDistance);
        m_Animator2.PlayMove();
    }

    public void Dash() {
        m_State = State.Dash;
        m_Animator2.SetMove(m_DashGroup.clips, m_RunGroup.moveDistance);
        m_Animator2.PlayMove();
        m_DashEffect.SetActive(true);
    }

    public void SetMove(string type) {
        if (type == null)
            m_RunGroup = GetRunGroup();
        else
            m_RunGroup = GetRunGroup(type);
        m_Animator2.SetMove(m_RunGroup.clips, m_RunGroup.moveDistance);
    }

    private RunData GetRunGroup(string type) {
        var group = new RunData(8);
        group[0] = GetData($"Run{type}_F");
        group[1] = GetData($"Run{type}_FL");
        group[2] = GetData($"Run{type}_L");
        group[3] = GetData($"Run{type}_BL");
        group[4] = GetData($"Run{type}_B");
        group[5] = GetData($"Run{type}_BR");
        group[6] = GetData($"Run{type}_R");
        group[7] = GetData($"Run{type}_FR");
        return group;
    }

    private RunData GetRunGroup() {
        var group = new RunData(8);
        group[0] = GetData("Run_F");
        group[1] = GetData("Run_FL");
        group[2] = GetData("Run_L");
        group[3] = GetData("Run_BL");
        group[4] = GetData("Run_B");
        group[5] = GetData("Run_BR");
        group[6] = GetData("Run_R");
        group[7] = GetData("Run_FR");
        return group;
    }

    private RunData GetDashGroup() {
        var group = new RunData(8);
        group[0] = GetData("Sprint");
        group[1] = GetData("Run_FL");
        group[2] = GetData("Run_L");
        group[3] = GetData("Run_BL");
        group[4] = GetData("Run_B");
        group[5] = GetData("Run_BR");
        group[6] = GetData("Run_R");
        group[7] = GetData("Run_FR");
        return group;
    }

    public void SetMode(string name) {
        m_Mode = name;
        m_IdleClip = GetData("Idle").clip;
        m_RunGroup = GetRunGroup();
        m_DashGroup = GetDashGroup();
        if (m_Animator2 != null) {
            m_Animator2.SetIdle(m_IdleClip);
            if (m_State == State.Dash)
                m_Animator2.SetMove(m_DashGroup.clips, m_DashGroup.moveDistance);
            else
                m_Animator2.SetMove(m_RunGroup.clips, m_RunGroup.moveDistance);
            if (m_Mode == "Common") {
                if (m_State == State.Idle)
                    m_Animator2.PlayIdle();
                m_Animator2.CloseBlend();
            }
        }
    }

    private CharacterAnimationDataClip GetData(string name) => BattleResources.GetAnimationDataClip(m_Posture, m_Mode, name);


    private class RunData {
        public AnimationClip[] clips;
        public float[] moveDistance;
        public RunData(int len) {
            clips = new AnimationClip[len];
            moveDistance = new float[len];
        }
        public CharacterAnimationDataClip this[int idx] {
            set {
                clips[idx] = value.clip;
                moveDistance[idx] = (float)value.offset;
            }
        }
    }

    private void SetMatirals() {
        foreach (Transform node in GetComponentsInChildren<Transform>()) {
            var renderer = node.GetComponent<Renderer>();
            if (!renderer)
                continue;
            m_Materials.AddRange(renderer.materials);
        }
        ResetRim();
        ResetSplash();
    }

    private void ResetRim() {
        foreach (var mat in m_Materials) {
            mat.SetInt("_UseRim", 0);
            mat.SetColor("_RimColor", new Color(1, 1, 1, 1));
            mat.SetTexture("_RimColorTex", null);

            mat.SetFloat("_RimMainStrength", 0);
            mat.SetFloat("_RimEnableLighting", 0);
            mat.SetFloat("_RimShadowMask", 0);
            mat.SetInt("_RimBackfaceMask", 0);
            mat.SetFloat("_RimBlendMode", 1);

            mat.SetFloat("_RimDirStrength", 0);
            mat.SetFloat("_RimIndirBorder", 1);
            mat.SetFloat("_RimBlur", 1);
            mat.SetFloat("_RimFresnelPower", 2);
            mat.SetFloat("_RimVRParallaxStrength", 0);
        }
    }

    private void ResetSplash() {
        foreach (var mat in m_Materials) {
            mat.SetInt("_UseEmission2nd", 0);
            mat.SetTexture("_Emission2ndMap", null);
            mat.SetColor("_Emission2ndColor", new Color(1, 1, 1, 1));
        }
    }

    private void UpdateRim() {
        if (m_DefRimColorLerp > 0) {
            m_DefRimColorLerp -= Time.deltaTime * 10;

            var o = m_DefRimCurColor;
            var n = Color.Lerp(m_DefRimTarColor, m_DefRimLasColor, m_DefRimColorLerp);

            if (o.a == 0 && n.a > 0)
                foreach (var mat in m_Materials)
                    mat.SetInt("_UseRim", 1);
            if (o.a > 0 && n.a == 0)
                foreach (var mat in m_Materials)
                    mat.SetInt("_UseRim", 0);
            foreach (var mat in m_Materials)
                mat.SetColor("_RimColor", n);

            m_DefRimCurColor = n;
        }
    }

    public void SetDefRim(int lv) {
        switch (lv) {
            case 1:
                m_DefRimTarColor = new Color(0, 0, 2, 1);
                m_DefRimLasColor = m_DefRimColorLerp == 0 ? new Color(0, 0, 2, 0) : m_DefRimCurColor;
                break;
            case 2:
                m_DefRimTarColor = new Color(2, 0, 2, 1);
                m_DefRimLasColor = m_DefRimColorLerp == 0 ? new Color(2, 0, 2, 0) : m_DefRimCurColor;
                break;
            case 3:
                m_DefRimTarColor = new Color(2, 1, 0, 1);
                m_DefRimLasColor = m_DefRimColorLerp == 0 ? new Color(2, 1, 0, 0) : m_DefRimCurColor;
                break;
            default:
                m_DefRimTarColor = m_DefRimLasColor = m_DefRimCurColor;
                m_DefRimTarColor.a = 0;
                break;
        }
        m_DefRimColorLerp = 1;
    }

    public void UpdateSplash() {
        if (m_SplashLerp > 0) {
            m_SplashLerp -= Time.deltaTime * 10;
            var o = m_SplashValue;
            var n = Mathf.Lerp(0, .3f, m_SplashLerp);

            if (o == 0 && n > 0)
                foreach (var mat in m_Materials)
                    mat.SetInt("_UseEmission2nd", 1);
            if (o > 0 && n == 0)
                foreach (var mat in m_Materials)
                    mat.SetInt("_UseEmission2nd", 0);

            foreach (var mat in m_Materials)
                mat.SetColor("_Emission2ndColor", new Color(1, 1, 1, n));

            m_SplashValue = n;
        }
    }

    public void SetSplash() {
        m_SplashLerp = 1;
    }

    public void Impulse(TSVector power) {
        var s = gameObject.AddComponent<CinemachineImpulseSource>();
        s.m_ImpulseDefinition.m_ImpulseType = CinemachineImpulseDefinition.ImpulseTypes.Uniform;
        s.m_ImpulseDefinition.m_ImpulseShape = CinemachineImpulseDefinition.ImpulseShapes.Bump;
        s.m_DefaultVelocity = power.ToVector();
        s.GenerateImpulseWithForce(1);
        Destroy(s, 2);
    }

    public void SetSpeed(float value) => m_Speed = value;

    public void UpdateEffect(List<FX> fxs) {
        for (int i = 0; i < m_FXRecord.Count; i++) {
            var fx = m_FXRecord[i];
            if (!fxs.Contains(fx)) {
                Destroy(m_FXDirectors[fx].gameObject);
                m_FXDirectors.Remove(fx);
                m_FXRecord.RemoveAt(i--);
            }
        }
        for (int i = 0; i < fxs.Count; i++) {
            var fx = fxs[i];
            if (!m_FXDirectors.ContainsKey(fx)) {
                GameObject go = Instantiate(fx.target, transform);
                var tr = go.transform;
                tr.localPosition = fx.position;
                tr.localEulerAngles = fx.rotation;
                tr.localScale = fx.scale;
                if (!fx.follow)
                    tr.SetParent(null, true);
                var director = go.GetComponent<PlayableDirector>();
                director.timeUpdateMode = DirectorUpdateMode.Manual;
                m_FXDirectors.Add(fx, director);
                m_FXRecord.Add(fx);
            }
        }
    }
}