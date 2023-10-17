using System;
using System.Collections.Generic;
using BDFramework.UFlux;
using BDFramework.UFlux.Collections;
using BDFramework.UFlux.View.Props;
using MR.Net.Proto.Battle;
using TMPro;
using UnityEngine.UI;

namespace UI.GameUI
{
  
    // public static class MsgDataPool
    // {
    //     public static Dictionary<Type, UIMsgData> msgDataDic = new Dictionary<Type, UIMsgData>();
    //     public static T GetMsgData<T>() where T : UIMsgData, new()
    //     {
    //         var type = typeof(T);
    //         msgDataDic.TryGetValue(type, out UIMsgData value);
    //         if (value == null)
    //         {
    //             value = new T();
    //             msgDataDic[type] = value;
    //         }
    //         return value as T;
    //     }
    // }


    public class UIMsg_Loading : UIMsgData
    {
        public UIMsg_Loading()
        {
        }

        public WinEnum winEnum;
        public bool isOpen;
        public Action callBack;
    }
}