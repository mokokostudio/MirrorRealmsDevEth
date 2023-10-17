using BDFramework.UFlux;
using BDFramework.UFlux.Collections;
using BDFramework.UFlux.View.Props;
using MR.Net.Frame;
using MR.Net.Proto.Battle;
using System.Collections.Generic;
using System.Linq;
using TMPro;
using UnityEngine;
using UnityEngine.UI;

public class SubWeaponItem : APropsBase
{
    [ComponentValueBind("Select", typeof(TransformHelper), nameof(TransformHelper.ShowHideTransForm))]
    public bool select;
    [ComponentValueBind("NftIcon", typeof(TransformHelper), nameof(TransformHelper.ShowHideTransForm))]
    public bool nftIcon;
    [ComponentValueBind("Icon", typeof(Image), nameof(Image.sprite))]
    public string weaponIcon;
    [ComponentValueBind("Icon", typeof(TransformHelper), nameof(TransformHelper.ShowHideTransForm))]
    public bool weapon;

    public WeaponPB weaponData;
}

public class Prop_Weapon : APropsBase
{
    [ComponentValueBind("ScrollRect/Viewport/Content", typeof(UFluxBindLogic), nameof(UFluxBindLogic.BindChildren))]
    public PropsList<SubWeaponItem> weaponList = new PropsList<SubWeaponItem>();
}

public class SubWeaponWindow : AWindow<Prop_Weapon>
{
    public SubWeaponWindow(string path) : base(path)
    {

    }

    public SubWeaponWindow(Transform transform) : base(transform)
    {
    }

    public static string[] qualityImg = { "White", "Green", "Blue", "Purple", "Gold" };
    public WeaponPB curWeapon;
    public WeaponPB equipWeapon;
    public int curClickIndex;
    [TransformPath("Name")]
    public TMP_Text weaponName;
    [TransformPath("Des")]
    public TMP_Text weaponDes;
    [TransformPath("Icon")]
    public Image weaponIcon;
    [TransformPath("SwitchBtn")]
    public Button swtichBtn;
    [TransformPath("NFT")]
    public Transform nftGO;
    [TransformPath("NFT/Progress")]
    public TMP_Text nftProgress;
    [TransformPath("NFT/NftDes")]
    public TMP_Text nftDes;
    [TransformPath("NFT/Slider")]
    public Slider nftSlider;
    [TransformPath("NFT/Slider/Quality")]
    public Image nftQuality;
    [TransformPath("NFT/Prop1")]
    public TMP_Text nftProp1;
    [TransformPath("NFT/Prop2")]
    public TMP_Text nftProp2;
    [TransformPath("NFT/UID")]
    public TMP_Text nftUID;
    [ButtonOnclick("CloseBtn")]
    public void OnClose()
    {
        this.Close();
    }
    public override void Init()
    {
        base.Init();
        for (int i = 0; i < 9; i++)
        {
            var item = new SubWeaponItem();
            item.select = false;
            item.weapon = false;
            item.nftIcon = false;
            Props.weaponList.Add(item);
        }
        Props.SetPropertyChange(nameof(Props.weaponList));
        CommitProps();
    }
    public override void Open(UIMsgData uiMsg = null)
    {
        var sData = uiMsg as SubWeaponData;
        curWeapon = sData.data;
        equipWeapon = sData.data;
        curClickIndex = sData.index;
        var wType = Config.Equips.Weapon[curWeapon.ConfigID].Type;
        List<WeaponPB> weaponList = new List<WeaponPB>();
        for (int i = 0; i < PlayerData.Weapons.Length; i++) {
            var item = PlayerData.Weapons[i];
            if (Config.Equips.Weapon[item.ConfigID].Type == wType)
                weaponList.Add(item);
        }
        for (int i = 0; i < 9; i++) {
            var item = Props.weaponList.Get(i);
            if (i < weaponList.Count) {
                item.weaponIcon = ResourcePath.WeaponPath + Config.Equips.Weapon[weaponList[i].ConfigID].Icon;
                item.select = weaponList[i] == curWeapon;
                item.weapon = true;
                item.nftIcon = weaponList[i].ID > 0;
                item.weaponData = weaponList[i];
                item.Transform.GetComponent<Button>().onClick.AddListener(() => { OnClickWeapon(item.weaponData); });
            } else {
                item.select = false;
                item.weapon = false;
                item.nftIcon = false;
                item.Transform.GetComponent<Button>().onClick.RemoveAllListeners();
            }
        }
        Props.SetPropertyChange(nameof(Props.weaponList));
        CommitProps();
        ChoseWeapon();
        base.Open(uiMsg);
    }

    public void OnClickWeapon(WeaponPB weaponData)
    {
        curWeapon = weaponData;
        for (int i = 0; i < Props.weaponList.Count; i++)
        {
            var item = Props.weaponList.Get(i);
            item.select = item.weaponData == curWeapon;
        }
        Props.SetPropertyChange(nameof(Props.weaponList));
        CommitProps();
        ChoseWeapon();
    }
    public void ChoseWeapon()
    {
        this.swtichBtn.interactable = equipWeapon != curWeapon;
        var config = Config.Equips.Weapon[curWeapon.ConfigID];
        weaponName.text = config.Name;
        weaponIcon.sprite = UFluxUtils.LoadSprite($"{ResourcePath.WeaponPath + config.Icon}_m");
        weaponDes.text = config.Info;
        if (config.Quality > Config.Global.QualityType.N)
        {
            nftGO.gameObject.SetActive(true);
            weaponDes.gameObject.SetActive(false);
            nftDes.text = config.Info;
            nftSlider.value = curWeapon.Quality / 100f;
            nftProgress.text = curWeapon.Quality.ToString();
            nftQuality.sprite = UFluxUtils.LoadSprite("Assets/Arts/UI/Tips/Progress_" + qualityImg[curWeapon.Quality / 20]);
            nftProp1.text = curWeapon.Prop1;
            nftProp2.text = curWeapon.Prop2;
            nftUID.text = $"#{curWeapon.ID:000000}";
        }
        else {
            nftGO.gameObject.SetActive(value: false);
            weaponDes.gameObject.SetActive(true);
        }
    }
    [ButtonOnclick("SwitchBtn")]
    public void OnSwitch()
    {
        var weaponIndex = PlayerData.Weapons.ToList().IndexOf(curWeapon);
        var tcp = Net.GetTcp();
        tcp.Send(new ChangeEquipC2S { LocationIndex = curClickIndex, WeaponIndex = weaponIndex }, (ChangeEquipS2C msg) =>
        {
            if (msg.Code == CodePBType.Success)
            {
                PlayerData.Equiped = msg.Equiped;
                var uimsg = new Window_Main.UIMsg_Main();
                uimsg.weaponID = curWeapon.ConfigID;
                UIManager.Inst.SendMessage(WinEnum.Win_Main, uimsg);
            }
            else
            {
                Debug.Log(msg.Code.ToString());
            }
        });
        this.Close();
    }

    public override void Close()
    {
        base.Close();
    }
    public class SubWeaponData : UIMsgData
    {
        public int index;
        public WeaponPB data;
    }
}

