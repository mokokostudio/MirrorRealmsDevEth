using Cinemachine;
using MR.Battle;
using System;
using System.Collections;
using System.Collections.Generic;
using UnityEngine;

[AddComponentMenu("")] // Don't display in add component menu
[SaveDuringPlay]
public class CinemachineBattleFree : CinemachineComponentBase {
    public override bool IsValid => enabled;
    public override CinemachineCore.Stage Stage => CinemachineCore.Stage.Aim;
    public float x;

    private Quaternion m_Target;
    private bool m_First = true;

    public void Reset() {
        m_First = true;
    }

    public override void MutateCameraState(ref CameraState curState, float deltaTime) {
        if (Battle.Instance == null)
            return;
        if (m_First) {
            m_Target = Quaternion.Euler(x, Battle.Instance.CameraDir.AsFloat() * Mathf.Rad2Deg, 0);
            m_First = false;
        } else {
            var target = Quaternion.Euler(x, Battle.Instance.CameraDir.AsFloat() * Mathf.Rad2Deg, 0);
            m_Target = Quaternion.Lerp(m_Target, target, deltaTime * 10);
        }
        curState.RawOrientation = m_Target;
    }
}
