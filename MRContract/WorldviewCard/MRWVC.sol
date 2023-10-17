// SPDX-License-Identifier: Unlicensed
pragma solidity >=0.8.4 <0.9.0;

import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/token/ERC721/extensions/ERC721Enumerable.sol";
import "@openzeppelin/contracts/token/ERC721/extensions/ERC721URIStorage.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/Counters.sol";
import "@openzeppelin/contracts/utils/math/SafeMath.sol";
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";

enum Level {Diamond, Gold, Silver, Bronze}
enum CardType {King, Queen, Rook, Bishop, Knight, Pawn}
enum BlindBoxType {OrdinaryBlindBox, LegendaryBlindBox}

contract MyYlNft is ERC721Enumerable, ERC721URIStorage, Ownable, ReentrancyGuard {
    using Counters for Counters.Counter;
    Counters.Counter private _tokenIds;
    using SafeMath for uint;

    struct YlWorldCard {
        uint256 tokenId; 
        string name; 
        uint256 price; 
        string uri; 
        uint16 ratio; 
        CardType cardType; 
        Level level; 
    }

    mapping(uint => YlWorldCard) private tokenOfYlWorldCard; //nft tokenid relate worldview card
    mapping(address => mapping(BlindBoxType => uint)) private ownedBoxCount; //buyer and amount

    uint private ordinaryBlindBoxPrice = 1000000000000000; //normal price
    uint private legendaryBlindBoxPrice = 2000000000000000; //legend price

    uint16 private totalOrdinaryBlindBoxQuantity = 3376; 
    uint16 private totalLegendaryBlindBoxQuantity = 1600; 
    uint16 private onSellOrdinaryBlindBoxQuantity = 0; 
    uint16 private onSellLegendaryBlindBoxQuantity = 0;


    uint16 private worldCardsQuantity = 4976; 
    uint16 private surplusWorldCardsQuantity = 4976; 
    uint16 private ordinaryBoxOfWorldCardsQuantity = 3376; 
    uint16 private legendaryBoxOfWorldCardsQuantity = 1600; 
    YlWorldCard[] private ordinaryBoxOfWorldCards; 
    YlWorldCard[] private legendaryBoxOfWorldCards; 

    address[] private ordinaryBlindBoxBuyer; 
    address[] private legendaryBlindBoxBuyer; 

    constructor () ERC721("YlNft", "MRWVC") {}

    function buyOrdinaryBlindBox(uint16 quantity) external payable {
        require(onSellOrdinaryBlindBoxQuantity >= quantity, "not enough onSellOrdinaryBlindBox");
        require(msg.value >= ordinaryBlindBoxPrice.mul(quantity), "not enough balance");
        buyBindBox(msg.sender, BlindBoxType.OrdinaryBlindBox , quantity);
    }

    function buyLegendaryBlindBox(uint16 quantity) external payable {
        require(onSellLegendaryBlindBoxQuantity >= quantity, "not enough onSellLegendaryBlindBox");
        require(msg.value >= legendaryBlindBoxPrice.mul(quantity), "not enough balance");
        buyBindBox(msg.sender, BlindBoxType.LegendaryBlindBox , quantity);
    }

    function giveBlindBox(address to, BlindBoxType boxType, uint16 quantity) public onlyOwner {
        buyBindBox(to, boxType, quantity);
    }

    function buyBindBox(address to, BlindBoxType boxType, uint16 quantity) internal {
        if (boxType == BlindBoxType.OrdinaryBlindBox) {
            require(onSellOrdinaryBlindBoxQuantity >= quantity, "not enough onSellOrdinaryBlindBox");
            onSellOrdinaryBlindBoxQuantity -= quantity;
            for (uint i = 1; i <= quantity; i++) {
                ordinaryBlindBoxBuyer.push(to);
            }
        } else if (boxType == BlindBoxType.LegendaryBlindBox) {
            require(onSellLegendaryBlindBoxQuantity >= quantity, "not enough onSellLegendaryBlindBox");
            onSellLegendaryBlindBoxQuantity -= quantity;
            for (uint i = 1; i <= quantity; i++) {
                legendaryBlindBoxBuyer.push(to);
            }
        }
        ownedBoxCount[to][boxType] += quantity;
    }


    function mintNftCardsToBindBox(
        string[] memory names,
        uint256[] memory prices,
        string[] memory uris,
        uint16[] memory ratios,
        CardType[] memory cardTypes,
        Level[] memory levels,
        BlindBoxType[] memory boxTypes,
        uint16[] memory batchSizes
    ) 
        external 
        onlyOwner 
    {
        require(names.length == prices.length, "size must be equation");
        require(prices.length == uris.length, "size must be equation");
        require(uris.length == ratios.length, "size must be equation");
        require(ratios.length == cardTypes.length, "size must be equation");
        require(cardTypes.length == levels.length, "size must be equation");
        require(levels.length == boxTypes.length, "size must be equation");
        require(boxTypes.length == batchSizes.length, "size must be equation");
        require(getSum(batchSizes) <= surplusWorldCardsQuantity, "over mintable quantity");
        require(getSum(batchSizes) > 0, "batchSize must greater than 0");
        require(batchSizes.length <= 4976, "batchSizes length too large");
        for (uint16 i = 0; i < batchSizes.length; i++) {
            if (boxTypes[i] == BlindBoxType.OrdinaryBlindBox) {
                require(levels[i] != Level.Diamond, "OrdinaryBlindBox can not put Diamond card");
                for (uint j = 0; j < batchSizes[i]; j++) {
                    YlWorldCard memory ordinary = mintYlWorldCardNft(names[i], prices[i], uris[i], ratios[i], cardTypes[i], levels[i]);
                    ordinaryBoxOfWorldCards.push(ordinary);
                }
            } else if (boxTypes[i] == BlindBoxType.LegendaryBlindBox) {
                for (uint j = 0; j < batchSizes[i]; j++) {
                    YlWorldCard memory legendary = mintYlWorldCardNft(names[i], prices[i], uris[i], ratios[i], cardTypes[i], levels[i]);
                    legendaryBoxOfWorldCards.push(legendary);
                }
            }
            surplusWorldCardsQuantity--;
        }
    }

    function mintYlWorldCardNft(
        string memory name,
        uint256 price,
        string memory uri,
        uint16 ratio,
        CardType cardType,
        Level level
    ) 
        internal
        returns (YlWorldCard memory)
    {
        uint tokenId = mintNft(msg.sender, uri);
        tokenOfYlWorldCard[tokenId] = YlWorldCard(tokenId, name, price, uri, ratio, cardType, level);
        return tokenOfYlWorldCard[tokenId];
    } 

    function openBlindBox() public onlyOwner {
        openOrdinaryBlindBox();
        openLegendaryBlindBox();
    }

    function openOrdinaryBlindBox() public onlyOwner {
        require(ordinaryBoxOfWorldCards.length >= ordinaryBlindBoxBuyer.length, "ordinary cards not enough");
        address[] memory ordinaryBuyer = ordinaryBlindBoxBuyer;
        for (uint i = 0; i < ordinaryBuyer.length; i++) {
            uint random = getRandom(ordinaryBoxOfWorldCards.length);
            YlWorldCard memory card = ordinaryBoxOfWorldCards[random];
            ordinaryBoxOfWorldCards[random] = ordinaryBoxOfWorldCards[ordinaryBoxOfWorldCards.length - 1];
            ordinaryBoxOfWorldCards.pop();
            ownedBoxCount[ordinaryBuyer[i]][BlindBoxType.OrdinaryBlindBox] -= 1;
            safeTransferFrom(msg.sender, ordinaryBuyer[i], card.tokenId);
        }
        delete ordinaryBlindBoxBuyer;
    }

    function openLegendaryBlindBox() public onlyOwner {
        require(legendaryBoxOfWorldCards.length >= legendaryBlindBoxBuyer.length, "legendary cards not enough");
        address[] memory legendaryBuyer = legendaryBlindBoxBuyer;
        for (uint i = 0; i < legendaryBuyer.length; i++) {
            uint random = getRandom(legendaryBoxOfWorldCards.length);
            YlWorldCard memory card = legendaryBoxOfWorldCards[random];
            legendaryBoxOfWorldCards[random] = legendaryBoxOfWorldCards[legendaryBoxOfWorldCards.length - 1];
            legendaryBoxOfWorldCards.pop();
            ownedBoxCount[legendaryBuyer[i]][BlindBoxType.LegendaryBlindBox] -= 1;
            safeTransferFrom(msg.sender, legendaryBuyer[i], card.tokenId);
        }
        delete legendaryBlindBoxBuyer;
    }

    function setBlindBoxPrice(BlindBoxType boxType, uint price) public onlyOwner {
        if (boxType == BlindBoxType.OrdinaryBlindBox) {
            ordinaryBlindBoxPrice = price;
        } else if (boxType == BlindBoxType.LegendaryBlindBox) {
            legendaryBlindBoxPrice = price;
        }
    }

    function setBlindBoxOnSell(BlindBoxType boxType, uint16 quantity) public onlyOwner {
        if (boxType == BlindBoxType.OrdinaryBlindBox) {
            require(totalOrdinaryBlindBoxQuantity >= quantity, "not enough ordinaryBlindBox");
            onSellOrdinaryBlindBoxQuantity += quantity;
            totalOrdinaryBlindBoxQuantity -= quantity;
        } else if (boxType == BlindBoxType.LegendaryBlindBox) {
            require(totalLegendaryBlindBoxQuantity >= quantity, "not enough legendaryBlindBox");
            onSellLegendaryBlindBoxQuantity += quantity;
            totalLegendaryBlindBoxQuantity -= quantity;
        }
    }

    function mintNft(address to, string memory uri) internal returns(uint256) {
        _tokenIds.increment();
        uint256 tokenId = _tokenIds.current();
        _mint(to, tokenId);
        _setTokenURI(tokenId, uri);
        return tokenId;
    }

    function getOwnedBoxCount(address o) public view returns(uint[] memory) {
        uint[] memory count = new uint[](2);
        count[0] = ownedBoxCount[o][BlindBoxType.OrdinaryBlindBox];
        count[1] = ownedBoxCount[o][BlindBoxType.LegendaryBlindBox];
        return count;
    }

    function getNftCards(address o) public view returns(YlWorldCard[] memory)  {
        require(balanceOf(o) > 0, "no nftCards");
        YlWorldCard[] memory cards = new YlWorldCard[](balanceOf(o));
        for (uint i = 0; i < balanceOf(o); i++)  {
            YlWorldCard memory card = tokenOfYlWorldCard[tokenOfOwnerByIndex(o, i)];
            cards[i] = card;
        }
        return cards;
    }

    function getOrdinaryCards() public view returns(YlWorldCard[] memory) {
        return ordinaryBoxOfWorldCards;
    }

    function getLegendaryCards() public view returns(YlWorldCard[] memory) {
        return legendaryBoxOfWorldCards;
    }

    function getRandom(uint number) public view returns(uint) {
        return uint(keccak256(abi.encodePacked(block.timestamp, block.prevrandao, 
            msg.sender))) % number;
    }

    function getSum(uint16[] memory batchSize) internal pure returns (uint) {
        uint sum = 0; 
        for (uint16 i = 0; i < batchSize.length; i++) {
            sum = sum.add(batchSize[i]);
        }
        return sum;
    }

    function _beforeTokenTransfer(address from, address to, uint256 tokenId, uint256 batchSize)
        internal
        override(ERC721, ERC721Enumerable)
    {
        super._beforeTokenTransfer(from, to, tokenId, batchSize);
    }

    function _burn(uint256 tokenId) internal override(ERC721, ERC721URIStorage) {
        super._burn(tokenId);

        if (bytes(tokenOfYlWorldCard[tokenId].name).length != 0) {
            delete tokenOfYlWorldCard[tokenId];
        }
    }

    function tokenURI(uint256 tokenId)
        public
        view
        override(ERC721, ERC721URIStorage)
        returns (string memory)
    {
        return super.tokenURI(tokenId);
    }

    function supportsInterface(bytes4 interfaceId)
        public
        view
        override(ERC721, ERC721Enumerable)
        returns (bool)
    {
        return super.supportsInterface(interfaceId);
    }
}