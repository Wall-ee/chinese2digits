from  decimal import Decimal
import re


CHINESE_CHAR_LIST = ['幺','零', '一', '二', '两', '三', '四', '五', '六', '七', '八', '九', '十', '百', '千', '万', '亿']
CHINESE_SIGN_LIST = ['负','正','-','+']
CHINESE_CONNECTING_SIGN_LIST = ['.','点','·']
CHINESE_PER_COUNTING_STRING_LIST = ['百分之','千分之','万分之']
CHINESE_PER_COUNTING_SEG = '分之'
CHINESE_PURE_NUMBER_LIST = ['幺', '一', '二', '两', '三', '四', '五', '六', '七', '八', '九', '十','零']

CHINESE_SIGN_DICT = {'负':'-','正':'+','-':'-','+':'+'}
CHINESE_PER_COUNTING_DICT = {'百分之':'%','千分之':'‰','万分之':'‱'}
CHINESE_CONNECTING_SIGN_DICT = {'.':'.','点':'.','·':'.'}
CHINESE_COUNTING_STRING = {'十':10, '百':100, '千':1000, '万':10000, '亿':100000000}
CHINESE_PURE_COUNTING_UNIT_LIST = ['十','百','千','万','亿']

TRADITIONAl_CONVERT_DICT = {'壹': '一', '贰': '二', '叁': '三', '肆': '四', '伍': '五', '陆': '六', '柒': '七',
                           '捌': '八', '玖': '九'}
SPECIAL_TRADITIONAl_COUNTING_UNIT_CHAR_DICT = {'拾': '十', '佰': '百', '仟':'千', '萬':'万', '億':'亿'}

SPECIAL_NUMBER_CHAR_DICT = {'两':'二','俩':'二'}

"""
中文转阿拉伯数字
"""
common_used_ch_numerals = {'幺':1,'零':0, '一':1, '二':2, '两':2, '三':3, '四':4, '五':5, '六':6, '七':7, '八':8, '九':9, '十':10, '百':100, '千':1000, '万':10000, '亿':100000000}

"""
阿拉伯数字转中文
"""
digits_char_ch_dict = {'0':'零','1':'一','2':'二','3':'三','4':'四','5':'五','6':'六','7':'七','8':'八','9':'九','%':'百分之','‰':'千分之','‱':'万分之','.':'点'}


#以百分号作为大逻辑区分。 是否以百分号作为新的数字切割逻辑 所以同一套切割逻辑要有  或关系   有百分之结尾 或者  没有百分之结尾
# takingChineseNumberRERules = re.compile(r'(?:(?:[正负]){0,1}(?:(?:[一二三四五六七八九十千万亿兆幺零百]+(?:点[一二三四五六七八九万亿兆幺零]+){0,1})|(?:点[一二三四五六七八九万亿兆幺零]+)))'
#                                         r'(?:(?:分之)(?:[正负]){0,1}(?:(?:[一二三四五六七八九十千万亿兆幺零百]+(?:点[一二三四五六七八九万亿兆幺零]+){0,1})|'
#                                         r'(?:点[一二三四五六七八九万亿兆幺零]+))){0,1}')

takingChineseDigitsMixRERules = re.compile(r'(?:(?:分之){0,1}(?:\+|\-){0,1}[正负]{0,1})'
                                            r'(?:(?:(?:\d+(?:\.\d+){0,1}(?:[\%\‰\‱]){0,1}|\.\d+(?:[\%\‰\‱]){0,1}){0,1}'
                                            r'(?:(?:(?:[一二三四五六七八九十千万亿兆幺零百]+(?:点[一二三四五六七八九万亿兆幺零]+){0,1})|(?:点[一二三四五六七八九万亿兆幺零]+))))'
                                            r'|(?:(?:\d+(?:\.\d+){0,1}(?:[\%\‰\‱]){0,1}|\.\d+(?:[\%\‰\‱]){0,1})'
                                            r'(?:(?:(?:[一二三四五六七八九十千万亿兆幺零百]+(?:点[一二三四五六七八九万亿兆幺零]+){0,1})|(?:点[一二三四五六七八九万亿兆幺零]+))){0,1}))')

PURE_DIGITS_RE = re.compile('[0-9]')




DIGITS_CHAR_LIST = ['0','1', '2', '3', '4', '5', '6', '7', '8', '9']
DIGITS_SIGN_LIST = ['-','+']
DIGITS_CONNECTING_SIGN_LIST = ['.']
DIGITS_PER_COUNTING_STRING_LIST = ['%','‰','‱']
takingDigitsRERule = re.compile(r'(?:(?:\+|\-){0,1}\d+(?:\.\d+){0,1}(?:[\%\‰\‱]){0,1}|(?:\+|\-){0,1}\.\d+(?:[\%\‰\‱]){0,1})')

def coreCHToDigits(chineseChars):
    total = 0
    tempVal = '' #用以记录临时是否建议数字拼接的字符串 例如 三零万 的三零
    countingUnit = 1              #表示单位：个十百千,用以计算单位相乘 例如八百万 百万是相乘的方法，但是如果万前面有 了一千八百万 这种，千和百不能相乘，要相加...
    countingUnitFromString = [1]   #原始字符串提取的单位应该是一个list  在计算的时候，新的单位应该是本次取得的数字乘以已经发现的最大单位，例如 4千三百五十万， 等于 4000万+300万+50万
    for i in range(len(chineseChars) - 1, -1, -1):
        val = common_used_ch_numerals.get(chineseChars[i])
        if val >= 10 and i == 0:  #应对 十三 十四 十*之类，说明为十以上的数字，看是不是十三这种
            #说明循环到了第一位 也就是最后一个循环 看看是不是单位开头
            #取最近一次的单位
            if val > countingUnit:  #如果val大于 contingUnit 说明 是以一个更大的单位开头 例如 十三 千二这种
                countingUnit = val   #赋值新的计数单位
                total = total + val    #总值等于  全部值加上新的单位 类似于13 这种
                countingUnitFromString.append(val)
            else:
                countingUnitFromString.append(val)
                # 计算用的单位是最新的单位乘以字符串中最大的原始单位  为了计算四百万这种
                # countingUnit = countingUnit * val
                countingUnit = max(countingUnitFromString) * val
                #total =total + r * x
        elif val >= 10:
            if val > countingUnit:
                countingUnit = val
                countingUnitFromString.append(val)
            else:
                countingUnitFromString.append(val)
                # 计算用的单位是最新的单位乘以字符串中最大的原始单位 为了计算四百万这种
                # countingUnit = countingUnit * val
                countingUnit = max(countingUnitFromString) * val
        else:
            if i > 0 :
                #如果下一个不是单位 则本次也是拼接
                if common_used_ch_numerals.get(chineseChars[i-1]) <10:
                    tempVal = str(val) + tempVal
                else:
                    #说明已经有大于10的单位插入 要数学计算了
                    #先拼接再计算
                    #如果取值不大于10 说明是0-9 则继续取值 直到取到最近一个大于10 的单位   应对这种30万20千 这样子
                    total = total + countingUnit * int(str(val) + tempVal)
                    #计算后 把临时字符串置位空
                    tempVal = ''
            else:
                #那就是无论如何要收尾了
                #如果counting unit 等于1  说明所有字符串都是直接拼接的，不用计算，不然会丢失前半部分的零
                if countingUnit == 1:
                    tempVal = str(val) + tempVal
                else:
                    total = total + countingUnit * int(str(val) + tempVal)

    #如果 total 为0  但是 countingUnit 不为0  说明结果是 十万这种  最终直接取结果 十万
    #如果countingUnit 大于10 说明他是就是 汉字零
    if total == 0:
        if countingUnit>10:
            total = str(countingUnit)
        else:
            if tempVal != "":
                total = tempVal
            else:
                total = str(total)
    else:
        total = str(total)
    return total
def chineseToDigits(chineseDigitsMixString,percentConvert = True,*args,**kwargs):
    #之前已经做过罗马数字变汉字 所以不存在罗马数字部分问题了
    """
    分之  分号切割  要注意
    """
    chineseCharsListByDiv = chineseDigitsMixString.split('分之')
    convertResultList = []
    for k in range(len(chineseCharsListByDiv)):
        tempChineseChars = chineseCharsListByDiv[k]

        # chineseChars = str(chineseChars)
        # tempChineseChars = chineseChars
        #kaka
        chineseCharsDotSplitList = []

        """
        看有没有符号
        """
        sign = ''
        for chars in tempChineseChars:
            if CHINESE_SIGN_DICT.get(chars) is not None:
                sign = CHINESE_SIGN_DICT.get(chars)
                tempChineseChars = tempChineseChars.replace(chars, '')
        """
        防止没有循环完成就替换 报错
        """
        chineseChars = tempChineseChars

        """
        小数点切割，看看是不是有小数点
        """
        for chars in list(CHINESE_CONNECTING_SIGN_DICT.keys()):
            if chars in chineseChars:
                chineseCharsDotSplitList = chineseChars.split(chars)

        if chineseCharsDotSplitList.__len__()==0:
            convertResult = coreCHToDigits(chineseChars)
        else:
            #如果小数点右侧有 单位 比如 2.55万  4.3百万 的处理方式
            #先把小数点右侧单位去掉
            tempCountString = ''
            for ii in range(len(chineseCharsDotSplitList[-1]) - 1, -1, -1):
                if chineseCharsDotSplitList[-1][ii] in CHINESE_PURE_COUNTING_UNIT_LIST:
                    tempCountString = chineseCharsDotSplitList[-1][ii] + tempCountString
                else:
                    chineseCharsDotSplitList[-1] = chineseCharsDotSplitList[-1][0:(ii+1)]
                    break
            if tempCountString != '':
                tempCountNum = Decimal(coreCHToDigits(tempCountString))
            else:
                tempCountNum = Decimal(1.0)
            if chineseCharsDotSplitList[0] == '':
                """
                .01234 这种开头  用0 补位
                """
                convertResult = '0.'+ coreCHToDigits(chineseCharsDotSplitList[1])
            else:
                """
                小数点右侧要注意，有可能是00开头
                """
                convertResult = coreCHToDigits(chineseCharsDotSplitList[0]) + '.' + coreCHToDigits(chineseCharsDotSplitList[1])

            convertResult = str(Decimal(convertResult) * tempCountNum)
        """
        如果 convertResult 是空字符串， 表示可能整体字符串是 负百分之10 这种  或者 -百分之10
        """
        if convertResult =='':
            convertResult = '1'

        convertResult = sign + convertResult

        # #处理小数点右边的0
        # if '.' in convertResult:
        #     convertResult = convertResult.rstrip('0')
        #     if convertResult.endswith('.'):
        #         convertResult = convertResult.rstrip('.')
        convertResultList.append(convertResult)
    if len(convertResultList)>1:
        #是否转换分号及百分比
        if percentConvert == True:
            finalTotal = str(Decimal(convertResultList[1])/Decimal(convertResultList[0]))
        else:
            if convertResultList[0] == '100':
                finalTotal = convertResultList[1] + '%'
            elif convertResultList[0] == '1000':
                finalTotal = convertResultList[1] + '‰'
            elif convertResultList[0] == '10000':
                finalTotal = convertResultList[1] + '‱'
            else:
                finalTotal = convertResultList[1]+'/' + convertResultList[0]
    else:
        finalTotal = convertResultList[0]
    # 处理小数点右边的0
    if '.' in finalTotal:
        finalTotal = finalTotal.rstrip('0')
        finalTotal = finalTotal.rstrip('.')
        # if finalTotal.endswith('.'):
        #     finalTotal = finalTotal.rstrip('.')
    return finalTotal


def chineseToDigitsHighTolerance(chineseDigitsMixString,percentConvert = True, skipError=False, errorChar=[], errorMsg=[]):
    if skipError:
        try:
            total = chineseToDigits(chineseDigitsMixString,percentConvert = percentConvert)
        except Exception as e:
            #返回类型不能是none 是空字符串
            total = ''
            errorChar.append(chineseDigitsMixString)
            errorMsg.append(str(e))
    else:
        total = chineseToDigits(chineseDigitsMixString,percentConvert = percentConvert)
    return total


def checkChineseNumberReasonable(chNumber):
    if chNumber.__len__()>0:
        #由于在上个检查点 已经把阿拉伯数字转为中文 因此不用检查阿拉伯数字部分
        """
        如果汉字长度大于0 则判断是不是 万  千  单字这种
        """
        for i in CHINESE_PURE_NUMBER_LIST:
            if i in chNumber:
                #只要有数字在字符串 就说明不是 千千万万这种只有单位的表述
                return True
    return False

"""
繁体简体转换 及  单位  特殊字符转换 两千变二千
"""
def traditionalTextConvertFunc(chString,traditionalConvertSwitch=True):
    chStringList = list(chString)
    stringLength = len(list(chStringList))

    if traditionalConvertSwitch == True:

        for i in range(stringLength):
            #繁体中文数字转简体中文数字
            if TRADITIONAl_CONVERT_DICT.get(chStringList[i],'') != '':
                chStringList[i] = TRADITIONAl_CONVERT_DICT.get(chStringList[i],'')
    if stringLength > 1:
        #检查繁体单体转换
        for i in range(stringLength):
            #如果 前后有 pure 汉字数字 则转换单位为简体
            if SPECIAL_TRADITIONAl_COUNTING_UNIT_CHAR_DICT.get(chStringList[i],'') != '':
                # 如果前后有单纯的数字 则进行单位转换
                if i == 0:
                    if chStringList[i+1] in CHINESE_PURE_NUMBER_LIST:
                        chStringList[i] = SPECIAL_TRADITIONAl_COUNTING_UNIT_CHAR_DICT.get(chStringList[i], '')
                elif i == stringLength-1:
                    if chStringList[i-1] in CHINESE_PURE_NUMBER_LIST:
                        chStringList[i] = SPECIAL_TRADITIONAl_COUNTING_UNIT_CHAR_DICT.get(chStringList[i], '')
                else:
                    if chStringList[i-1] in CHINESE_PURE_NUMBER_LIST or \
                            chStringList[i+1] in CHINESE_PURE_NUMBER_LIST :
                        chStringList[i] = SPECIAL_TRADITIONAl_COUNTING_UNIT_CHAR_DICT.get(chStringList[i], '')
            #特殊变换 俩变二
            if SPECIAL_NUMBER_CHAR_DICT.get(chStringList[i], '') != '':
                # 如果前后有单位 则进行转换
                if i == 0:
                    if chStringList[i+1] in CHINESE_PURE_COUNTING_UNIT_LIST:
                        chStringList[i] = SPECIAL_NUMBER_CHAR_DICT.get(chStringList[i], '')
                elif i == stringLength-1:
                    if chStringList[i-1] in CHINESE_PURE_COUNTING_UNIT_LIST:
                        chStringList[i] = SPECIAL_NUMBER_CHAR_DICT.get(chStringList[i], '')
                else:
                    if chStringList[i-1] in CHINESE_PURE_COUNTING_UNIT_LIST or \
                            chStringList[i+1] in CHINESE_PURE_COUNTING_UNIT_LIST :
                        chStringList[i] = SPECIAL_NUMBER_CHAR_DICT.get(chStringList[i], '')
    return ''.join(chStringList)

"""
标准表述转换  三千二 变成 三千二百 三千十二变成 三千零一十二 四万十五变成四万零十五
"""
def standardChNumberConvert(chNumberString):
    chNumberStringList = list(chNumberString)

    #大于2的长度字符串才有检测和补位的必要
    if len(chNumberStringList) > 2:
        #十位补一：
        try:
            tenNumberIndex = chNumberStringList.index('十')
            if tenNumberIndex == 0:
                chNumberStringList.insert(tenNumberIndex, '一')
            else:
                # 如果没有左边计数数字 插入1
                if chNumberStringList[tenNumberIndex - 1] not in CHINESE_PURE_NUMBER_LIST:
                    chNumberStringList.insert(tenNumberIndex, '一')
        except:
            pass

        #差位补零
        #逻辑 如果最后一个单位 不是十结尾 而是百以上 则数字后面补一个比最后一个出现的单位小一级的单位
        #从倒数第二位开始看,且必须是倒数第二位就是单位的才符合条件
        try:
            lastCountingUnit = CHINESE_PURE_COUNTING_UNIT_LIST.index(chNumberStringList[len(chNumberStringList)-2])
            # 如果最末位的是百开头
            if lastCountingUnit >= 1:
                # 则字符串最后拼接一个比最后一个单位小一位的单位 例如四万三 变成四万三千

                # 如果最后一位结束的是亿 则补千万
                if lastCountingUnit == 4:
                    chNumberStringList.append('千万')
                else:
                    chNumberStringList.append(CHINESE_PURE_COUNTING_UNIT_LIST[lastCountingUnit - 1])
        except:
            pass
    #检查是否是 万三  千四点五这种表述 百三百四
    perCountSwitch = 0
    if len(chNumberStringList) >1:
        if chNumberStringList[0] in ['千','万','百']:
            for i in range(1,len(chNumberStringList)):
                #其余位数都是纯数字 才能执行
                if chNumberStringList[i] in CHINESE_PURE_NUMBER_LIST:
                    perCountSwitch = 1
                else:
                    perCountSwitch = 0
                    #y有一个不是数字 直接退出循环
                    break
    if perCountSwitch == 1:
       chNumberStringList = chNumberStringList[:1]+['分','之']+chNumberStringList[1:]
    return ''.join(chNumberStringList)


def checkNumberSeg(chineseNumberList,originText):
    newChineseNumberList = []
    #用来控制是否前一个已经合并过  防止多重合并
    tempPreText = ''
    tempMixedString = ''
    segLen = len(chineseNumberList) 
    if segLen >0:
        #加入唯一的一个 或者第一个
        if chineseNumberList[0][:2] in CHINESE_PER_COUNTING_SEG:
            #如果以分之开头 记录本次 防止后面要用 是否出现连续的 分之
            tempPreText = chineseNumberList[0]
            newChineseNumberList.append(chineseNumberList[0][2:])
        else:
            newChineseNumberList.append(chineseNumberList[0])

        if len(chineseNumberList)>1:
            for i in range(1,segLen):
                #判断本字符是不是以  分之  开头  
                if chineseNumberList[i][:2] in CHINESE_PER_COUNTING_SEG:
                    #如果是以 分之 开头 那么检查他和他见面的汉子数字是不是连续的 即 是否在原始字符串出现
                    tempMixedString = chineseNumberList[i-1] + chineseNumberList[i]
                    if tempMixedString in originText:
                        #如果连续的上一个字段是以分之开头的  本字段又以分之开头  
                        if tempPreText != '':
                            #检查上一个字段的末尾是不是 以 百 十 万 的单位结尾
                            if tempPreText[-1] in CHINESE_PURE_COUNTING_UNIT_LIST:
                                #先把上一个记录进去的最后一位去掉
                                newChineseNumberList[-1] = newChineseNumberList[-1][:-1]
                                #如果结果是确定的，那么本次的字段应当加上上一个字段的最后一个字
                                newChineseNumberList.append(tempPreText[-1] + chineseNumberList[i])
                            else:
                                #如果上一个字段不是以单位结尾  同时他又是以分之开头，那么 本次把分之去掉
                                newChineseNumberList.append(chineseNumberList[i][2:])
                        else:
                            #上一个字段不以分之开头，那么把两个字段合并记录
                            if newChineseNumberList.__len__()>0:
                                newChineseNumberList[-1] = tempMixedString
                            else:
                                newChineseNumberList.append(tempMixedString)
                    else:
                        #说明前一个数字 和本数字不是连续的
                        #本数字去掉分之二字
                        newChineseNumberList.append(chineseNumberList[i][2:])

                    #记录以 分之 开头的字段  用以下一个汉字字段判别
                    tempPreText = chineseNumberList[i]
                else:
                    #不是  分之 开头 那么把本数字加入序列
                    newChineseNumberList.append(chineseNumberList[i])
                    #记录把不是 分之 开头的字段  临时变量记为空
                    tempPreText = ''
    return newChineseNumberList

def checkSignSeg(chineseNumberList):
    newChineseNumberList = []
    tempSign = ''
    for i in range(len(chineseNumberList)):
        #新字符串 需要加上上一个字符串 最后1位的判断结果
        newChNumberString = tempSign  + chineseNumberList[i]
        lastString = newChNumberString[-1:]
        #如果最后1位是正负号 那么本字符去掉最后1位  下一个数字加上最后3位
        if lastString in CHINESE_SIGN_LIST:
            tempSign = lastString
            #如果最后1位 是  那么截掉最后1位
            newChNumberString = newChNumberString[:-1]
        else:
            tempSign = ''
        newChineseNumberList.append(newChNumberString)
    return newChineseNumberList

def digitsToCHChars(mixedStringList):
    resultList = []
    for mixedString in mixedStringList:
        if mixedString.startswith('.'):
            mixedString = '0'+mixedString
        for key in digits_char_ch_dict.keys():
            if key in mixedString:
                # 应当记录下来有转换，然后再操作  在核心函数里 通过小数点判断是否应该强制  
                mixedString = mixedString.replace(key,digits_char_ch_dict.get(key))
                #应当是只要有百分号 就挪到前面 阿拉伯数字没有四百分之的说法
                #防止这种 3%万 这种问题
                for k in CHINESE_PER_COUNTING_STRING_LIST:
                    if k in mixedString:
                        temp = k + mixedString.replace(k,'')
                        mixedString = temp

        resultList.append(mixedString)
    return resultList



def takeChineseNumberFromString(chText,percentConvert = True,traditionalConvert= True,digitsNumberSwitch= False,verbose=False,*args,**kwargs):
    """
    :param chText: chinese string
    :param percentConvert: convert percent simple. Default is True.  3% will be 0.03 in the result
    :param traditionalConvert: Switch to convert the Traditional Chinese character to Simplified chinese
    :param digitsNumberSwitch: Switch to convert the take pure digits number
    :return: Dict like result. 'inputText',replacedText','CHNumberStringList':CHNumberStringList,'digitsStringList'
    """

    """
    是否只提取数字
    """
    if digitsNumberSwitch is True:
        return takeDigitsNumberFromString(chText,percentConvert=percentConvert)


    """
    简体转换开关
    """
    # originText = chText

    convertedCHString = traditionalTextConvertFunc(chText,traditionalConvert)

    """
    字符串 汉字数字字符串切割提取
    正则表达式方法
    """
    CHNumberStringListTemp = takingChineseDigitsMixRERules.findall(convertedCHString)
    #检查是不是  分之 切割不完整问题
    CHNumberStringListTemp = checkNumberSeg(CHNumberStringListTemp,convertedCHString)

    #检查末位是不是正负号
    CHNumberStringListTemp = checkSignSeg(CHNumberStringListTemp)

    #备份一个原始的提取，后期处结果的时候显示用
    OriginCHNumberTake = CHNumberStringListTemp.copy()

    #将阿拉伯数字变成汉字  不然合理性检查 以及后期 如果不是300万这种乘法  而是 四分之345  这种 就出错了
    CHNumberStringListTemp = digitsToCHChars(CHNumberStringListTemp)


    #检查合理性 是否是单纯的单位  等
    CHNumberStringList= []
    OriginCHNumberForOutput = []
    for i in range(len(CHNumberStringListTemp)):
        tempText = CHNumberStringListTemp[i]
        if checkChineseNumberReasonable(tempText):
            #如果合理  则添加进被转换列表
            CHNumberStringList.append(tempText)
            #则添加把原始提取的添加进来
            OriginCHNumberForOutput.append(OriginCHNumberTake[i])
    #TODO 检查是否 时间格式 五点四十  七点一刻

    """
    进行标准汉字字符串转换 例如 二千二  转换成二千零二
    """
    CHNumberStringListTemp = list(map(lambda x:standardChNumberConvert(x),CHNumberStringList))

    """
    将中文转换为数字
    """
    digitsStringList = []
    replacedText = convertedCHString
    errorCharList = []
    errorMsgList = []
    if CHNumberStringListTemp.__len__()>0:
        # digitsStringList = list(map(lambda x:chineseToDigits(x,percentConvert=percentConvert),CHNumberStringList))
        # 用标准清洗后的字符串进行转换
        # digitsStringList = list(
        #     map(lambda x: chineseToDigits(x,percentConvert=percentConvert), CHNumberStringListTemp))

        for kk in range(len(CHNumberStringListTemp)):
            # digitsStringList.append(chineseToDigits(CHNumberStringListTemp[kk], percentConvert=percentConvert))
            digitsStringList.append(chineseToDigitsHighTolerance(CHNumberStringListTemp[kk],percentConvert=percentConvert,skipError=verbose,errorChar=errorCharList,errorMsg=errorMsgList))
        # tupleToReplace = list(zip(OriginCHNumberForOutput,digitsStringList,list(map(len,OriginCHNumberForOutput))))
        tupleToReplace = [ (d,c,i) for d,c,i in zip(OriginCHNumberForOutput,digitsStringList,list(map(len,OriginCHNumberForOutput))) if c !='']


        """
        按照提取出的中文数字字符串长短排序，然后替换。防止百分之二十八 ，二十八，这样的先把短的替换完了的情况
        """
        tupleToReplace = sorted(tupleToReplace, key=lambda x: -x[2])
        for i in range(tupleToReplace.__len__()):
            # if tupleToReplace[i][0] is None:
            #     continue
            replacedText = replacedText.replace(tupleToReplace[i][0],tupleToReplace[i][1])


    # finalResult = {
    #     'inputText':originText,
    #     'replacedText':replacedText,
    #     'CHNumberStringList':OriginCHNumberForOutput,
    #     'digitsStringList':digitsStringList
    # }
    finalResult = {
        'inputText':chText,
        'replacedText':replacedText,
        'CHNumberStringList':OriginCHNumberForOutput,
        'digitsStringList':digitsStringList,
        'errorWordList': errorCharList,
        'errorMsgList': errorMsgList
    }
    return finalResult


def takeNumberFromString(chText,percentConvert = True,traditionalConvert= True,digitsNumberSwitch= False, verbose=False, *args,**kwargs):
    """
    :param chText: chinese string
    :param percentConvert: convert percent simple. Default is True.  3% will be 0.03 in the result
    :param traditionalConvert: Switch to convert the Traditional Chinese character to Simplified chinese
    :param digitsNumberSwitch: Switch to convert the take pure digits number
    :param verbose: if true, will return words that raised exception and catch the error
    :return: Dict like result. 'inputText',replacedText','CHNumberStringList':CHNumberStringList,'digitsStringList'
    """
    finalResult = takeChineseNumberFromString(chText,percentConvert = percentConvert,traditionalConvert= traditionalConvert,digitsNumberSwitch= digitsNumberSwitch,verbose=verbose)
    return finalResult



def takeDigitsNumberFromString(textToExtract,percentConvert = False):
    digitsNumberStringList = takingDigitsRERule.findall(textToExtract)
    """
    最后检查有没有百分号
    """
    digitsStringList = []
    replacedText = textToExtract
    if digitsNumberStringList.__len__()>0:
        # digitsStringList = list(map(lambda x:chineseToDigits(x,percentConvert=percentConvert),digitsNumberStringList))
        tupleToReplace = list(zip(digitsNumberStringList,digitsStringList,list(map(len,digitsNumberStringList))))


        """
        按照提取出的中文数字字符串长短排序，然后替换。防止百分之二十八 ，二十八，这样的先把短的替换完了的情况
        """
        tupleToReplace = sorted(tupleToReplace, key=lambda x: -x[2])
        for i in range(tupleToReplace.__len__()):
            replacedText = replacedText.replace(tupleToReplace[i][0],tupleToReplace[i][1])

    finalResult = {
        'inputText':textToExtract,
        'replacedText':replacedText,
        'digitsNumberStringList':digitsNumberStringList,
        'digitsStringList':digitsStringList
    }
    return finalResult

if __name__=='__main__':
    print(takeNumberFromString('3.1千万'))

    print(takeNumberFromString('拾'))

    print(takeNumberFromString('12.55万'))

    #混合提取
    print(takeNumberFromString('三零万二零千拉阿拉啦啦30万20千嚯嚯或百四嚯嚯嚯四百三十二分之2345啦啦啦啦',percentConvert=False))
    print(takeNumberFromString('百分之5负千分之15'))
    print(takeNumberFromString('啊啦啦啦300十万你好我20万.3%万你好啊300咯咯咯-.34%啦啦啦300万'))
    print(takeChineseNumberFromString('百分之四百三十二万分之四三千分之五今天天气不错三百四十点零零三四'))
    #将百分比转为小数
    print(takeDigitsNumberFromString('234%lalalal-%nidaye+2.34%',percentConvert=True))
    print(takeChineseNumberFromString('aaaa.3%万'))
    # 新测试用例  分数测试  四分之三 这种
    print(takeNumberFromString('十分之一'))
    print(takeChineseNumberFromString('四分之三啦啦五百分之二',percentConvert=False))
    print(takeChineseNumberFromString('4分之3负五分之6咿呀呀 四百分之16ooo千千万万'))
    print(takeNumberFromString('百分之五1234%'))
    print(takeNumberFromString('五百分之一',percentConvert=False))

    #正则引擎已经全部使用re2 规则 不再用pcre规则 防止出现递归炸弹
    print(takeChineseNumberFromString('百分之四百三十二万分之四三千分之五'))

    #测试繁体 简称  等
    print(takeNumberFromString('四千三'))
    print(takeChineseNumberFromString('伍亿柒仟万拾柒今天天气不错百分之三亿二百万五啦啦啦啦负百分之点二八你好啊三万二'))
    print(takeChineseNumberFromString('llalala万三威风威风千四五'))
    print(takeChineseNumberFromString('哥两好'))
    print(takeChineseNumberFromString('伍亿柒仟万拾柒百分之'))
    print(takeChineseNumberFromString('负百分之点二八你好啊百分之三五是不是点五零百分之负六十五点二八'))

    
