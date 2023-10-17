using System.Collections.Generic;
using System.Collections;
using System.Threading.Tasks;
using BDFramework.UFlux;
using Cysharp.Threading.Tasks;
using UnityEngine;
using UnityEngine.ResourceManagement.AsyncOperations;
using UnityEngine.UI;
using System;
using UnityEngine.AddressableAssets;
using BDFramework;
using BDFramework.ScreenView;
using UI.GameUI;

[UI((int) WinEnum.Win_Loading, "Window/Window_Loading")]
public class Window_Loding : AWindow
{
    public Window_Loding(string path) : base(path)
    {
    }

    [TransformPath("Adapter/Slider")] private Slider slider;
    //[TransformPath("Progress")] private Text proText;

    public override void Init()
    {
        base.Init();
    }
    [UIMessageListener]
    private void Message(UIMsg_Loading msg)
    {
        ProgressTask(msg).Forget();
    }
    private async UniTaskVoid ProgressTask(UIMsg_Loading msg)
    {
        var TaskList = UFluxUtils.TaskList;
        while (true)
        {
            float progress = 0;
            for (int i = 0; i < TaskList.Count; i++)
                progress += TaskList[i].PercentComplete;
            progress /= TaskList.Count;
            slider.value = Mathf.MoveTowards(slider.value, progress, Time.deltaTime);
            if (slider.value == 1)
            {
                if (msg.isOpen)
                    UIManager.Inst.ShowWindow(msg.winEnum);
                Close();
                msg.callBack?.Invoke();
                slider.value = 0;
                break;
            }

            await UniTask.Yield();
        }
    }
}