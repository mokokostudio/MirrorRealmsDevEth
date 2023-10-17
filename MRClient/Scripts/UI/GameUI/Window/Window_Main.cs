using BDFramework.UFlux;
using BDFramework.UFlux.Collections;
using BDFramework.UFlux.View.Props;
using Cysharp.Threading.Tasks;
using MR.Net.Proto.Battle;
using System;
using System.Collections;
using System.Collections.Generic;
using System.Linq;
using System.Threading;
using TMPro;
using UnityEngine;
using UnityEngine.Animations;
using UnityEngine.Networking;
using UnityEngine.UI;

public enum MainAniState
{
    Idle = 0,
    Arm = 1,    //肩部
    Hair = 2,   //头发
    Foot = 3,   //脚
    Crotch = 4, //裆
    Chest = 5,  //胸
    Waist = 6,  //腰
    Fight = 7,  //进入战场
    Abdomen = 8, //肚子

    Daggers = 10,//双匕
    SwordShield,//剑盾
    GreatSword,//巨剑
    Bow,//长弓
}
public class WeaponItem : APropsBase
{
    public ulong weaponID;
    [ComponentValueBind("Displaying", typeof(Button), nameof(Button.onClick))]
    public Action displaying;
    [ComponentValueBind("Exhibit", typeof(Button), nameof(Button.onClick))]
    public Action exhibit;
    [ComponentValueBind("Displaying", typeof(TransformHelper), nameof(TransformHelper.ShowHideTransForm))]
    public bool displayingGo;
    [ComponentValueBind("Exhibit", typeof(TransformHelper), nameof(TransformHelper.ShowHideTransForm))]
    public bool exhibitGo;
    [ComponentValueBind("IconBg/WeaponIcon", typeof(Button), nameof(Button.onClick))]
    public Action onSwitch;
    [ComponentValueBind("WeaponName", typeof(TMP_Text), nameof(TMP_Text.text))]
    public string weaponName;
    [ComponentValueBind("IconBg/WeaponIcon", typeof(Image), nameof(Image.sprite))]
    public string weaponIcon;
    [ComponentValueBind("IconBg/NftIcon", typeof(TransformHelper), nameof(TransformHelper.ShowHideTransForm))]
    public bool nftIcon;
}

public class Prop_Main : APropsBase
{
    [ComponentValueBind("Adapter/WeaponList", typeof(UFluxBindLogic), nameof(UFluxBindLogic.BindChildren))]
    public PropsList<WeaponItem> weaponList = new PropsList<WeaponItem>();
}

[UI((int)WinEnum.Win_Main, "Window/Window_Main")]
public class Window_Main : AWindow<Prop_Main>
{
    public static int RoomID;

    private int curWeaponIndex = -1;
    private int curClickWeapon = 0;
    private bool m_Init;

    [TransformPath("Adapter/Player/Eos")]
    public Transform m_Player;
    private List<GameObject> m_WeaponGo;
    private Animator m_Animator;
    private bool m_AniWeapon;
    public Window_Main(string path) : base(path)
    {
    }

    public override void Init()
    {
        base.Init();
        RegisterSubWindow(new SubRoomWindow(this.Transform.Find("Adapter/SubRoomPanel")));
        RegisterSubWindow(new SubWeaponWindow(this.Transform.Find("Adapter/SubWeaponPanel")));
        NameText.text = PlayerData.Name;
        Props.weaponList = new PropsList<WeaponItem>();
        m_WeaponGo = new List<GameObject>();
        for (int i = 0; i < PlayerData.Equiped.Length; i++)
        {
            var index = i;
            ulong wID = PlayerData.Weapons[PlayerData.Equiped[i]].ConfigID;
            var config = Config.Equips.Weapon[wID];
            var weaponItem = new WeaponItem();
            weaponItem.displayingGo = false;
            weaponItem.exhibitGo = true;
            weaponItem.weaponIcon = ResourcePath.WeaponPath + config.Icon;
            weaponItem.nftIcon = config.Quality > Config.Global.QualityType.N;
            weaponItem.weaponName = config.Name;
            weaponItem.weaponID = wID;
            weaponItem.exhibit = () => { OnExhiBit(index); };
            weaponItem.displaying = () => { OnDisplaying(index); };
            weaponItem.onSwitch = () => { OnShowWeapon(index); };
            Props.weaponList.Add(weaponItem);
        }
        Props.SetPropertyChange(nameof(Props.weaponList));
        CommitProps();
        m_Animator = m_Player.GetComponent<Animator>();
        m_Animator.SetInteger("AniState", (int)MainAniState.Idle);
        var rt = new RenderTexture(Screen.width, Screen.height, 0);
        Transform.Find("Bg/RawImage").GetComponent<RawImage>().texture = rt;
        Transform.Find("Adapter/Player/Camera").GetComponent<Camera>().targetTexture = rt;
        BGMManager.Play(BGMManager.BGMType.Main);
        //SetAnimator().Forget();
    }
    public override void Open(UIMsgData uiMsg = null)
    {
        base.Open();
        if (!m_Init)
        {
            var pTr = Transform.Find("Adapter/Player");
            pTr.position = Vector3.zero;
            pTr.localScale = Vector3.one / pTr.lossyScale.x;
            m_Init = true;
        }

    }
    [TransformPath("Adapter/UID/GameName")]
    public TMP_Text NameText;

    [ButtonOnclick("Adapter/CreatRoom")]
    private void OnReadyBattle()
    {
        RoomID = 0;
        UIManager.Inst.ShowWindow(WinEnum.Win_BattleReady);
    }

    [ButtonOnclick("Adapter/JoinRoom")]
    private void JoinBattle()
    {
        GetSubWindow<SubRoomWindow>().Open();
    }

    [ButtonOnclick("Adapter/PlayerIcon/HairBtn")]
    private void OnClickHair()
    {

        PlayAnimation2(() => { PlayAnimator(MainAniState.Hair); }, () => { PlayAnimator(MainAniState.Idle); });
    }
    [ButtonOnclick("Adapter/PlayerIcon/ChestBtn")]
    private void OnClickChest()
    {
        PlayAnimation2(() => { PlayAnimator(MainAniState.Chest); }, () => { PlayAnimator(MainAniState.Idle); });
    }

    [ButtonOnclick("Adapter/PlayerIcon/WaistBtn")]
    private void OnClickWaist()
    {
        PlayAnimation2(() => { PlayAnimator(MainAniState.Waist); }, () => { PlayAnimator(MainAniState.Idle); });
    }
    [ButtonOnclick("Adapter/PlayerIcon/FootBtn")]
    private void OnClickFoot()
    {

        PlayAnimation2(() => { PlayAnimator(MainAniState.Foot); }, () => { PlayAnimator(MainAniState.Idle); });
    }
    [ButtonOnclick("Adapter/PlayerIcon/ArmBtn")]
    private void OnClickArm()
    {
        PlayAnimation2(() => { PlayAnimator(MainAniState.Arm); }, () => { PlayAnimator(MainAniState.Idle); });
    }
    [ButtonOnclick("Adapter/PlayerIcon/ArmBtn2")]
    private void OnClickArm2()
    {
        PlayAnimation2(() => { PlayAnimator(MainAniState.Arm); }, () => { PlayAnimator(MainAniState.Idle); });
    }
    [ButtonOnclick("Adapter/PlayerIcon/CrotchBtn")]
    private void OnClickCrotch()
    {
        PlayAnimation2(() => { PlayAnimator(MainAniState.Crotch); }, () => { PlayAnimator(MainAniState.Idle); });
    }
    [ButtonOnclick("Adapter/PlayerIcon/AbdomenBtn")]
    private void OnClickAbdomen()
    {
        PlayAnimation2(() => { PlayAnimator(MainAniState.Abdomen); }, () => { PlayAnimator(MainAniState.Idle); });
    }
    private void PlayFight()
    {
        PlayAnimation2(() => { PlayAnimator(MainAniState.Fight); }, () => { PlayAnimator(MainAniState.Idle); });
    }


    private void OnChangeEquip(ChangeEquipS2C msg)
    {

    }
    public void PlayAnimator(MainAniState state)
    {

        if (state < MainAniState.Daggers && m_AniWeapon)
            return;
        // m_Player.localPosition = Vector3.zero;
        m_Animator.SetInteger("AniState", (int)state);
    }

    async UniTask SetAnimator()
    {

        while (true)
        {
            m_Player.localPosition = Vector3.zero;
            await UniTask.Yield();
        }
    }

    public void OnExhiBit(int index)
    {
        m_AniWeapon = true;
        if (curWeaponIndex != -1)
        {
            var item = Props.weaponList.Get(curWeaponIndex);
            item.displayingGo = false;
            item.exhibitGo = true;
        }
        var curItem = Props.weaponList.Get(index);
        curItem.displayingGo = true;
        curItem.exhibitGo = false;
        curWeaponIndex = index;
        curClickWeapon = index;
        Props.SetPropertyChange(nameof(Props.weaponList));
        CommitProps();
        OnWeapon();
    }

    public void OnDisplaying(int index)
    {
        m_AniWeapon = false;
        var curItem = Props.weaponList.Get(index);
        curItem.displayingGo = false;
        curItem.exhibitGo = true;
        curWeaponIndex = -1;
        Props.SetPropertyChange(nameof(Props.weaponList));
        CommitProps();
        for (int i = 0; i < m_WeaponGo.Count; i++)
        {
            GameObject.Destroy(m_WeaponGo[i].gameObject);
        }
        PlayAnimator(MainAniState.Idle);

    }

    public void OnWeapon()
    {
        if (curWeaponIndex != curClickWeapon)
            return;
        var wID = PlayerData.Weapons[PlayerData.Equiped[curClickWeapon]].ConfigID;
        var weapon = Config.Equips.Weapon[wID];
        var weaponType = Config.Equips.WeaponType[weapon.Type];
        int wType = (int)(weaponType.ID % 1000000000 / 1000) + 10;


        for (int i = 0; i < m_WeaponGo.Count; i++)
        {
            GameObject.Destroy(m_WeaponGo[i].gameObject);
        }
        m_WeaponGo.Clear();
        PlayAnimator((MainAniState)wType);
        for (int i = 0; i < weapon.Prefabs.Count; i++)
        {
            var wGo = GameObject.Instantiate(UFluxUtils.LoadAsync<GameObject>($"Weapon/{weapon.Prefabs[i]}.prefab"));
            var wp = m_Player.Find($"Weapons/{i + 1}");
            wp.localPosition = Vector3.zero;
            var pc = wp.GetComponent<ParentConstraint>();
            pc.weight = 0;
            wGo.transform.SetParent(wp, false);
            for (int j = 0; j < pc.sourceCount; j++)
            {
                var source = pc.GetSource(j);
                source.weight = 0;
                source.weight = weaponType.EquipPosition[i] == j ? 1 : 0;
                pc.SetSource(j, source);
            }
            pc.weight = 1;
            m_WeaponGo.Add(wGo);
        }

    }

    public void OnShowWeapon(int index)
    {
        SubWeaponWindow.SubWeaponData data = new SubWeaponWindow.SubWeaponData();
        data.index = index;
        data.data = PlayerData.Weapons[PlayerData.Equiped[index]];
        curClickWeapon = index;
        GetSubWindow<SubWeaponWindow>().Open(data);
    }

    public class UIMsg_Main : UIMsgData
    {
        public ulong weaponID;
    }
    [UIMessageListener]
    private void Message(UIMsg_Main msg)
    {
        var item = Props.weaponList.Get(curClickWeapon);
        var config = Config.Equips.Weapon[msg.weaponID];
        var weaponIndex = PlayerData.Weapons.ToList().FindIndex(m => m.ConfigID == msg.weaponID);
        item.weaponID = msg.weaponID;
        item.weaponName = config.Name;
        item.weaponIcon = ResourcePath.WeaponPath + config.Icon;
        item.nftIcon = config.Quality > Config.Global.QualityType.N;
        Props.SetPropertyChange(nameof(Props.weaponList));
        CommitProps();
        OnWeapon();
    }


    public void PlayAnimator(Animator animator, string clipName, Action startAct = null, Action endAct = null)
    {
        PlayAnimationItor(animator, clipName, startAct, endAct).Forget();
    }

    /// <summary>
    /// Animation动画播放迭代器
    /// </summary>
    /// <param name="animation">Animation组件</param>
    /// <param name="clipName">clip片段名</param>
    /// <param name="startAct">委托函数</param>
    /// <param name="endAct">委托函数</param>
    /// <returns></returns>
    private async UniTaskVoid PlayAnimationItor(Animator animator, string clipName, Action startAct, Action endAct)
    {
        startAct?.Invoke();
        AnimatorStateInfo animatorStateInfo = animator.GetCurrentAnimatorStateInfo(0);
        await new WaitForEndOfAnimr(animatorStateInfo, clipName);
        endAct?.Invoke();
    }
    private CancellationTokenSource m_AnimatorCancel;
    /// <summary>
    /// Animation动画播放迭代器
    /// </summary>
    /// <param name="animation">Animation组件</param>
    /// <param name="clipName">clip片段名</param>
    /// <param name="startAct">委托函数</param>
    /// <param name="endAct">委托函数</param>
    /// <returns></returns>
    async UniTask PlayAnimationItor2(Action startAct, Action endAct, CancellationToken token)
    {
        startAct?.Invoke();
        float timer = 0.01f;
        while (m_Animator.GetCurrentAnimatorClipInfo(0)[0].clip.name == "Girl_Common_Idle")
        {
            timer += Time.deltaTime;
            await UniTask.Yield();
        }

        try
        {
            await UniTask.Delay(TimeSpan.FromSeconds(m_Animator.GetCurrentAnimatorClipInfo(0)[0].clip.length - timer), cancellationToken: token);
        }
        catch (OperationCanceledException ex)
        {
            if (ex.CancellationToken == token)
            {
                return;
            }
        }
        endAct?.Invoke();
    }

    private void PlayAnimation2(Action startAct, Action endAct)
    {
        if (m_AnimatorCancel != null)
            m_AnimatorCancel.Cancel();
        m_AnimatorCancel = new CancellationTokenSource();
        PlayAnimationItor2(startAct, endAct, m_AnimatorCancel.Token).Forget();

    }
}

public class WaitForEndOfAnimr : IEnumerator
{
    AnimatorStateInfo m_animState;
    string m_clipName;
    public WaitForEndOfAnimr(AnimatorStateInfo animState, string clipName)
    {
        m_animState = animState;
        m_clipName = clipName;
    }
    public object Current
    {
        get
        {
            return null;
        }
    }
    public bool MoveNext()
    {
        return m_animState.IsName(m_clipName);
    }
    public void Reset()
    {
    }
}