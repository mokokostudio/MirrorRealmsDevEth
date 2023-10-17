using System;
using System.Collections;
using System.Collections.Generic;
using BDFramework.ScreenView;
using UnityEngine;
using BDFramework.UFlux;
using UnityEditor;
using UnityEngine.ResourceManagement.AsyncOperations;
using Cysharp.Threading.Tasks;
using UnityEngine.AddressableAssets;

[ScreenView((int) ScreenViewEnum.HotUpdate)]
public class ScreenView_Hot : IScreenView
{
    public bool IsLoad { get; private set; }
    public int Name { get; set; }

    public void BeginInit()
    {
        //一定要设置为true，否则当前是未加载状态
        IsLoad = true;
        Debug.Log($"进入热更");
        PlayerPrefs.SetInt("Main_", 1);
        UIManager.Inst.LoadWindow(WinEnum.Win_Loading,
            () =>
            {
                UIManager.Inst.ShowWindow(WinEnum.Win_Loading);
                ScreenViewManager.Inst.MainLayer.BeginNavTo(ScreenViewEnum.Login);
            });
    }

    public void BeginExit()
    {
    }
}