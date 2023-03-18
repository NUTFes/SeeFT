import 'package:seeft_mobile/configs/importer.dart';

class SignUpPage extends StatefulWidget {
  @override
  _SignUpPageState createState() => _SignUpPageState();
}

/// 利用者登録のページ
/// [_placeList]や[_ageList]などはAPIから取ってきたほうがいいかもしれない
class _SignUpPageState extends State<SignUpPage> {
  String _age = '';
  String _place = '';
  String _gender = '';
  List _selectedStores = [];

  List<String> _placeList = <String>[
    '',
    '新潟市北区',
    '新潟市東区',
    '新潟市中央区',
    '新潟市江南区',
    '新潟市秋葉区',
    '新潟市南区',
    '新潟市西区',
    '新潟市西蒲区',
    '新潟市外'
  ];

  List<String> _ageList = <String>[
    '',
    '～10代',
    '20代',
    '30代',
    '40代',
    '50代',
    '60代',
    '70代～'
  ];

  List<String> _storeList = [
    'ゴジバ',
    'エスメウ',
    'ああああ',
    'UNIKLO',
    'パン屋さん',
    'STARHACKS',
    'AGIGAS'
  ];

  /// setter
/*
  void _handleAge(String e) {
    setState(() {
      _age = e;
    });
  }

  void _handlePlace(String e) {
    setState(() {
      _place = e;
    });
  }

  void _handleGender(String e) {
    setState(() {
      _gender = e;
    });
  }
*/
  /// チェックボックスのチェックを可変にも対応できる。
  /// [selectedStores]に[store_id]があるかどうかで判断するためvalueではcontainsを使う。
  void _onStoreSelected(bool selected, storeID) {
    if (selected) {
      setState(() {
        _selectedStores.add(storeID);
      });
    } else {
      setState(() {
        _selectedStores.remove(storeID);
      });
    }
  }

  /// 画面描写
  @override
  Widget build(BuildContext context) {
    final Size size = MediaQuery.of(context).size;

    return Scaffold(
      appBar: AppBar(
        title: Text('Register User Page'),
      ),
      body: Container(
        padding: const EdgeInsets.all(40.0),
        child: SingleChildScrollView(
          child: Column(children: <Widget>[
            contentPlace(size),
            contentAge(size),
            contentGender(size),
            contentStore(size),
            RaisedButton(
              onPressed: _submission,
              child: Text('登録'),
            )
          ]),
        ),
      ),
    );
  }

  /// 地域選択
  Widget contentPlace(Size size) {
    return Row(
      children: [
        SizedBox(
          height: 50,
          width: size.width / 2 - 65,
          child: Align(
            alignment: Alignment.center,
            child: Text('お住いの地域'),
          ),
        ),
        SizedBox(width: 25),
        DropdownButton<String>(
          items: _placeList.map<DropdownMenuItem<String>>((String value) {
            return DropdownMenuItem<String>(
              value: value,
              child: Text(value),
            );
          }).toList(),
          onChanged: (e) {
            // ぬるぽちぇっく必要そう
            setState(() {
              _place = e!;
            });
          },
          value: _place,
        ),
      ],
    );
  }

  /// 年齢入力
  Widget contentAge(Size size) {
    return Row(
      children: [
        SizedBox(
          height: 50,
          width: size.width / 2 - 65,
          child: Align(
            alignment: Alignment.center,
            child: Text('年齢'),
          ),
        ),
        SizedBox(width: 45),
        DropdownButton<String>(
          items: _ageList.map<DropdownMenuItem<String>>((String value) {
            return DropdownMenuItem<String>(
              value: value,
              child: Text(value),
            );
          }).toList(),
          onChanged: (e) {
            // ぬるぽちぇっく必要そう
            setState(() {
              _age = e!;
            });
          },
          value: _age,
        ),
      ],
    );
  }

  /// 性別選択
  Widget contentGender(Size size) {
    return Row(
      children: [
        SizedBox(
          height: 50,
          width: size.width / 2 - 65,
          child: Align(
            alignment: Alignment.center,
            child: Text('性別'),
          ),
        ),
        Flexible(
          child: RadioListTile(
            activeColor: Colors.blue,
            title: Text('男'),
            value: '男',
            groupValue: _gender,
            onChanged: (e) {
              // ぬるぽちぇっく必要そう
              setState(() {
                _gender = e!.toString();
              });
            },
          ),
        ),
        Flexible(
          child: RadioListTile(
            activeColor: Colors.blue,
            title: Text('女'),
            value: '女',
            groupValue: _gender,
            onChanged: (e) {
              // ぬるぽちぇっく必要そう
              setState(() {
                _gender = e!.toString();
              });
            },
          ),
        ),
      ],
    );
  }

  /// 店舗選択
  Widget contentStore(Size size) {
    return Column(
      children: [
        Row(
          children: [
            SizedBox(
              width: size.width / 2 - 65,
              height: 50,
              child: Container(
                alignment: Alignment.center,
                child: Text('お気に入りの店舗'),
              ),
            ),
          ],
        ),
        SizedBox(height: 20),
      ],
    );
  }

  /// 登録ボタンをおしたときの処理
  /// ダイアログで入力したデータを確認できる
  Future _submission() async {
    logger.i('press register button.');
    var value = await showDialog(
      context: context,
      builder: (BuildContext context) => new AlertDialog(
        title: new Text('確認'),
        content: new Text(
            'place:$_place\nage:$_age\ngender:$_gender\nstore:$_selectedStores'),
        actions: <Widget>[
          new SimpleDialogOption(
            child: new Text('OK'),
            onPressed: () => Navigator.of(context).pop(),
          )
        ],
      ),
    );
    return value;
  }
}
