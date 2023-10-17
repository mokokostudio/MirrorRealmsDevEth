using BDFramework.UFlux;
using Cysharp.Threading.Tasks;
using System.Collections.Generic;
using System.Resources;
using UnityEngine;
using UnityEngine.AddressableAssets;
using UnityEngine.AddressableAssets.ResourceLocators;
using UnityEngine.InputSystem;
using UnityEngine.ResourceManagement.Util;

public static class BattleResources {

    private static Transform m_PoolRoot;
    public static Transform PoolRoot {
        get {
            if (!m_PoolRoot) {
                var go = new GameObject("Pool");
                Object.DontDestroyOnLoad(go);
                m_PoolRoot = go.transform;
            }
            return m_PoolRoot;
        }
    }

    private static Dictionary<string, List<GameObject>> m_Prefabs = new Dictionary<string, List<GameObject>>();
    private static Dictionary<string, Dictionary<string, Dictionary<string, AnimationData>>> m_AnimationDatas = new Dictionary<string, Dictionary<string, Dictionary<string, AnimationData>>>();

    public static void HidePool() {
        PoolRoot.gameObject.SetActive(false);
    }

    public static GameObject GetAvatar(string name) => GetPrefab($"Avatar_{name}");
    public static GameObject GetWeapon(string name) => GetPrefab($"Weapon_{name}");

    private static GameObject GetPrefab(string name) {
        if (m_Prefabs.ContainsKey(name)) {
            var list = m_Prefabs[name];
            if (list.Count > 1) {
                var go = list[0];
                go.transform.SetParent(null);
                return go;
            } else {
                var go = Object.Instantiate(list[0]);
                return go;
            }
        }
        return null;
    }

    public static CharacterAnimationDataClip GetAnimationDataClip(string posture, string mode, string name) {
        m_AnimationDatas.TryGetValue(posture, out var postureDic);
        if (postureDic == null)
            return null;
        postureDic.TryGetValue(mode, out var modeDic);
        if (modeDic == null)
            return null;
        modeDic.TryGetValue(name, out var data);
        if (data == null)
            return null;
        if (!data.dataClip) {
            var ao = data.reference.LoadAssetAsync<CharacterAnimationDataClip>();
            ao.WaitForCompletion();
            data.dataClip = ao.Result;
            data.reference.ReleaseAsset();
        }
        return data.dataClip;
    }

    public static void AsPreLoad() {
        AsPreLoadAvatars().Forget();
        AsPreLoadAnimations();
        AsPreLoadWeapons();
        PoolRoot.gameObject.SetActive(true);
    }

    private static async UniTaskVoid AsPreLoadAvatars() {
        await AsPreLoadAvatar("Eos");
        await AsPreLoadAvatar("EosBlack");
    }

    private static async UniTask AsPreLoadAvatar(string name) {
        var ao = UFluxUtils.LoadAsset<GameObject>($"Avatar/{name}.prefab");
        await ao;
        m_Prefabs[$"Avatar_{name}"] = new List<GameObject> { Object.Instantiate(ao.Result, PoolRoot) };
    }

    private static async UniTaskVoid AsPreLoadAnimation(string posture, string mode) {
        var ao = UFluxUtils.LoadAsset<CharacterAnimationDataGroup>($"Character/Animations/{posture}_{mode}/_Group.asset");
        await ao;
        m_AnimationDatas.TryGetValue(posture, out var postureDic);
        if (postureDic == null)
            m_AnimationDatas[posture] = postureDic = new Dictionary<string, Dictionary<string, AnimationData>>();

        postureDic.TryGetValue(mode, out var modeDic);
        if (modeDic == null)
            postureDic[mode] = modeDic = new Dictionary<string, AnimationData>();

        foreach (var data in ao.Result.dataClips)
            modeDic[data.name] = new AnimationData { reference = data.reference };
    }

    private static async UniTaskVoid AsPreLoadWeapon(string name) {
        var ao = UFluxUtils.LoadAsset<GameObject>($"Weapon/{name}.prefab");
        await ao;
        m_Prefabs[$"Weapon_{name}"] = new List<GameObject> { Object.Instantiate(ao.Result, PoolRoot) };
    }

    private static void AsPreLoadAnimations() {
        AsPreLoadAnimation("Girl", "Common").Forget();
        AsPreLoadAnimation("Girl", "SwordShield").Forget();
        //AsPreLoadAnimation("Girl", "TwinDaggers").Forget();
        AsPreLoadAnimation("Girl", "GreatSword").Forget();
        AsPreLoadAnimation("Girl", "Bow").Forget();
    }
    private static void AsPreLoadWeapons() {
        //AsPreLoadWeapon("Dagger").Forget();
        AsPreLoadWeapon("Shield").Forget();
        AsPreLoadWeapon("Sword").Forget();
        AsPreLoadWeapon("Bow").Forget();
        AsPreLoadWeapon("GreatSword").Forget();
        AsPreLoadWeapon("GreatSword2").Forget();
    }
    private class AnimationData {
        public CharacterAnimationDataClip dataClip;
        public AssetReference reference;
    }
}
