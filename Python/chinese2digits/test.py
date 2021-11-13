import unittest

import chinese2digits as c2d

class TestDict(unittest.TestCase):

    def test_mixedExtract(self):
        print('===========testing mixed string extract============')
        result = c2d.takeNumberFromString('三零万二零千拉阿拉啦啦30万20千嚯嚯或百四嚯嚯嚯四百三十二分之2345啦啦啦啦',percentConvert=False)
        self.assertEqual(result['replacedText'], '320000拉阿拉啦啦30000020000嚯嚯或4%嚯嚯嚯2345/432啦啦啦啦')
        self.assertEqual(result['CHNumberStringList'], ['三零万二零千', '30万', '20千', '百四', '四百三十二分之2345'])
        self.assertEqual(result['digitsStringList'], ['320000', '300000', '20000', '4%', '2345/432'])


        result = c2d.takeNumberFromString('百分之5负千分之15')
        self.assertEqual(result['replacedText'], '0.05-0.015')
        self.assertEqual(result['CHNumberStringList'], ['百分之5', '负千分之15'])
        self.assertEqual(result['digitsStringList'], ['0.05', '-0.015'])

        result = c2d.takeNumberFromString('拾')
        self.assertEqual(result['replacedText'], '拾')
        self.assertEqual(result['CHNumberStringList'], [])
        self.assertEqual(result['digitsStringList'], [])


        result = c2d.takeNumberFromString('零零零三四二啦啦啦啦12.550万啦啦啦啦啦零点零零三四二万')
        self.assertEqual(result['replacedText'], '000342啦啦啦啦125500啦啦啦啦啦34.2')
        self.assertEqual(result['CHNumberStringList'], ['零零零三四二','12.550万','零点零零三四二万'])
        self.assertEqual(result['digitsStringList'], ['000342', '125500','34.2'])




        result = c2d.takeNumberFromString('2.55万nonono3.1千万')
        self.assertEqual(result['replacedText'], '25500nonono31000000')
        self.assertEqual(result['CHNumberStringList'], ['2.55万','3.1千万'])
        self.assertEqual(result['digitsStringList'], ['25500','31000000'])

        result = c2d.takeNumberFromString('啊啦啦啦300十万你好我20万.3%万你好啊300咯咯咯-.34%啦啦啦300万')
        self.assertEqual(result['replacedText'], '啊啦啦啦30000000你好我20000030你好啊300咯咯咯-0.0034啦啦啦3000000')
        self.assertEqual(result['CHNumberStringList'], ['300十万', '20万', '.3%万', '300', '-.34%', '300万'])
        self.assertEqual(result['digitsStringList'], ['30000000', '200000', '30', '300', '-0.0034', '3000000'])

    def test_percentage_convert(self):
        print('===========testing percentage convert============')
        result = c2d.takeNumberFromString('aaaa.3%万啦啦啦啦0.03万')
        self.assertEqual(result['replacedText'], 'aaaa30啦啦啦啦300')
        self.assertEqual(result['CHNumberStringList'], ['.3%万','0.03万'])
        self.assertEqual(result['digitsStringList'], ['30','300'])

    def test_rational_number(self):
        print('===========testing rational number convert============')

        result = c2d.takeNumberFromString('十分之一')
        self.assertEqual(result['replacedText'], '0.1')
        self.assertEqual(result['CHNumberStringList'], ['十分之一'])
        self.assertEqual(result['digitsStringList'], ['0.1'])

        result = c2d.takeNumberFromString('四分之三啦啦五百分之二',percentConvert=False)
        self.assertEqual(result['replacedText'], '3/4啦啦2/500')
        self.assertEqual(result['CHNumberStringList'], ['四分之三', '五百分之二'])
        self.assertEqual(result['digitsStringList'], ['3/4', '2/500'])

        result = c2d.takeNumberFromString('4分之3负五分之6咿呀呀 四百分之16ooo千千万万')
        self.assertEqual(result['replacedText'], '0.75-1.2咿呀呀 4.16ooo千千万万')
        self.assertEqual(result['CHNumberStringList'], ['4分之3', '负五分之6', '四百分之16'])
        self.assertEqual(result['digitsStringList'], ['0.75', '-1.2', '4.16'])

        result = c2d.takeNumberFromString('百分之四百三十二万分之四三千分之五今天天气不错三百四十点零零三四')
        self.assertEqual(result['replacedText'], '4.320.00430.005今天天气不错340.0034')
        self.assertEqual(result['CHNumberStringList'], ['百分之四百三十二', '万分之四三', '千分之五','三百四十点零零三四'])
        self.assertEqual(result['digitsStringList'],  ['4.32', '0.0043', '0.005','340.0034'])

    def test_complex_convert(self):
        print('===========testing complex convert============')
        result = c2d.takeNumberFromString('四千三')
        self.assertEqual(result['replacedText'], '4300')
        self.assertEqual(result['CHNumberStringList'], ['四千三'])
        self.assertEqual(result['digitsStringList'], ['4300'])

        result = c2d.takeNumberFromString('伍亿柒仟万拾柒今天天气不错百分之三亿二百万五啦啦啦啦负百分之点二八你好啊三万二')
        self.assertEqual(result['replacedText'], '570000017今天天气不错3020050啦啦啦啦-0.0028你好啊32000')
        self.assertEqual(result['CHNumberStringList'],['五亿七千万十七', '百分之三亿二百万五', '负百分之点二八', '三万二'])
        self.assertEqual(result['digitsStringList'],  ['570000017', '3020050', '-0.0028', '32000'])

        result = c2d.takeNumberFromString('llalala万三威风威风千四五')
        self.assertEqual(result['replacedText'], 'llalala0.0003威风威风0.045')
        self.assertEqual(result['CHNumberStringList'], ['万三', '千四五'])
        self.assertEqual(result['digitsStringList'], ['0.0003', '0.045'])

        result = c2d.takeNumberFromString('伍亿柒仟万拾柒百分之')
        self.assertEqual(result['replacedText'], '570001700分之')
        self.assertEqual(result['CHNumberStringList'], ['五亿七千万十七百'])
        self.assertEqual(result['digitsStringList'], ['570001700'])

        result = c2d.takeNumberFromString('负百分之点二八你好啊百分之三五是不是点五零百分之负六十五点二八')
        self.assertEqual(result['replacedText'], '-0.0028你好啊0.35是不是0.5-0.6528')
        self.assertEqual(result['CHNumberStringList'], ['负百分之点二八', '百分之三五', '点五零', '百分之负六十五点二八'])
        self.assertEqual(result['digitsStringList'], ['-0.0028', '0.35', '0.5', '-0.6528'])

if __name__ == '__main__':
    unittest.main()