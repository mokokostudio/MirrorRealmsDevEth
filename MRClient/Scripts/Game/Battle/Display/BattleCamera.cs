using Cinemachine;
using System.Collections;
using System.Collections.Generic;
using UnityEngine;
namespace MR.Battle {
    public class BattleCamera : MonoBehaviour {
        public static BattleCamera Instance { get; private set; }
        public CinemachineVirtualCamera freeCam;
        public Transform m_WatchingTarget;

        private void Awake() {
            Instance = this;
        }

        private void Update() {
            if (Battle.Instance != null && Battle.Instance.CameraPlayer != null) {
                var go = Battle.Instance.CameraPlayer.Unit.GetComponentData<UnitDisplayCD>().Go;
                if (go) {
                    var target = go.transform;
                    if (target != m_WatchingTarget) {
                        freeCam.Follow = target.transform;
                        freeCam.UpdateCameraState(Vector3.up, 1000);
                    }
                }
            }
        }

        public void SetFollow(Transform target) {
            freeCam.Follow = target;
        }
    }
}