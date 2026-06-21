import 'package:seeft_mobile/configs/importer.dart';

//1ページ目とその表示
Future<void> openRescueDialog(BuildContext context,) async {
  return showDialog<void>(
    context: context,
    barrierDismissible: false, // user must tap button!
    barrierColor: Colors.black.withValues(alpha: 0.5),
    builder: (BuildContext context) {
      int rescuePageIndex = 0;
      bool manualChecked = false;
      bool poepleChecked = false;
      bool confirmError1 = false;

      bool taskChecked = false;
      bool locateChecked = false;
      bool contentChecked = false;
      bool confirmError2 = false;
      TextEditingController rescueContent = TextEditingController(); 
      TextEditingController rescueNumberofPeople = TextEditingController(); 
      String? rescuePostLevel;
      String? rescuePostLeve2;

      bool isChecked = false;


      return StatefulBuilder(
        builder: (context, setState){
          return AlertDialog(
            //各ページのタイトル
            title:Text(
              rescuePageIndex == 0 ? 'レスキュー要請1/2':
              rescuePageIndex == 1 ? 'レスキュー要請1/2':
              rescuePageIndex == 2 ? '最終確認':
              '要請完了'
            ),
            content:SizedBox(
              height:
              (rescuePageIndex == 0) ? MediaQuery.of(context).size.height*0.2 :
              (rescuePageIndex == 1) ? MediaQuery.of(context).size.height*0.5 :
              MediaQuery.of(context).size.height*0.1,
              child:IndexedStack(
              index:rescuePageIndex,
              children:[
                //1ページ目のcontent
                SingleChildScrollView(
                  child:ListBody(
                    children:<Widget>[
                      Text('確認', style:Theme.of(context).textTheme.titleLarge),
                      Text('本部に応援を要請する前に、以下のことを行ってください'),
                      Material(
                        color:Colors.transparent,
                        child:InkWell(
                          splashColor:Colors.transparent,
                          highlightColor:Colors.transparent,
                          onTap:(){
                                setState((){
                                  manualChecked = !manualChecked;
                                });
                              },
                          child:Row(
                            children:[
                              Checkbox(
                                value: manualChecked,
                                onChanged: (bool? value) {
                                  setState(() { // 🔹 `setState()` で UI を更新
                                    manualChecked = value!;
                                  });
                                },
                              ),
                              Text('マニュアルを確認した')
                            ],
                              ),
                          ),
                      ),
                      Material(
                        color:Colors.transparent,
                        child:InkWell(
                          splashColor:Colors.transparent,
                          highlightColor:Colors.transparent,
                          onTap:(){
                                setState((){
                                  poepleChecked= !poepleChecked;
                                });
                              },
                        child:Row(
                          children:[
                            Checkbox(
                            value: poepleChecked,
                            onChanged: (bool? value) {
                              setState(() { // 🔹 `setState()` で UI を更新
                                poepleChecked = value!;
                              });
                            },
                            ),
                            Text('近くの先輩に相談をした')

                          ]
                        )
                        ),
                      ),
                      if (confirmError1)
                        Text('確認をしてください',style:TextStyle(color:Colors.red))
                    ]
                  )
                ),
                //2ページ目のcontent
                SingleChildScrollView(
                  child:ListBody(
                      children:<Widget>[
                        Text('タスク', style:Theme.of(context).textTheme.titleLarge),
                        Text('問題の発生したタスクを選択してください'),
                        Autocomplete<String>(
                          optionsBuilder: (TextEditingValue textEditingValue) {
                            List<String> tasklist = ['駐車場誘導', '受付（講義棟前）', '受付（プール前）', 'その他'];

                            // 🔍 ユーザーの入力に基づいてリストをフィルタリング
                            return tasklist.where((String item) {
                              return item.toLowerCase().contains(textEditingValue.text.toLowerCase());
                            }).toList();
                          },
                          onSelected: (String selecttask) {
                            taskChecked = !taskChecked;
                          },
                          fieldViewBuilder: (context, textEditingController, focusNode, onFieldSubmitted) {
                            return TextField(
                              controller: textEditingController,
                              focusNode: focusNode,
                              decoration: InputDecoration(
                                labelText: "タスク",
                                prefixIcon: Icon(Icons.search), // 🔍 検索アイコンを追加
                                border: OutlineInputBorder(),
                              ),
                            );
                          },
                        ), 
                        if(taskChecked == false && confirmError2 == true)
                          Text('必要事項を入力してください',style:TextStyle(color:Colors.red)),
                        SizedBox(height: 20),
                        Text('場所', style:Theme.of(context).textTheme.titleLarge),
                        Text('現在の場所を教えてください'),
                        Autocomplete<String>(
                          optionsBuilder: (TextEditingValue textEditingValue) {
                            List<String> locatelist= ['体育館','プール前','講義棟前','駐車場入口','電気棟駐車場'];
                            return locatelist.where((String item){
                              return item.toLowerCase().contains(textEditingValue.text.toLowerCase());
                            }).toList();
                          },
                          onSelected: (String selectlocate) {
                            locateChecked = !locateChecked;
                          },
                          fieldViewBuilder: (context, textEditingController, focusNode, onFieldSubmitted) {
                            return TextField(
                              controller: textEditingController,
                              focusNode: focusNode,
                              decoration: InputDecoration(
                                labelText: "場所",
                                prefixIcon: Icon(Icons.search), // 🔍 検索アイコンを追加
                                border: OutlineInputBorder(),
                              ),
                            );
                          }
                        ),
                        if(locateChecked == false && confirmError2 == true)
                          Text('必要事項を入力してください',style:TextStyle(color:Colors.red)),
                        Text('内容', style:Theme.of(context).textTheme.titleLarge),
                        Text('どんな問題が発生したのか、なるべく詳しく記入してください。'),
                        TextField(
                          controller: rescueContent,
                          maxLines:null,
                          decoration:InputDecoration(
                            border: OutlineInputBorder(),
                            hintText:'Label'
                          )
                        ),
                        Text('必要人数',style:Theme.of(context).textTheme.titleLarge),
                        Text('必要な人数を入力してください(半角数字)'),
                        TextField(
                          controller: rescueNumberofPeople,
                          keyboardType: TextInputType.number,
                          inputFormatters:[FilteringTextInputFormatter.digitsOnly],
                          maxLines:1,
                          decoration:InputDecoration(
                            border:OutlineInputBorder(),
                            hintText:'Input'
                          )
                        ),
                        Text('要求レベル',style:Theme.of(context).textTheme.titleLarge),
                        Text('どの程度の役職が必要か選択してください'),
                        ListTile(
                          title:Text('部門長レベル'),
                          leading:Radio(
                            groupValue: rescuePostLevel,
                            value:'A',
                            onChanged:(String? value){
                              setState((){
                                rescuePostLevel = value;
                              });
                            },
                          ),
                          onTap:(){
                            setState((){
                              rescuePostLevel = 'A';
                            });
                          },
                        ),
                        ListTile(
                          title:Text('局長レベル'),
                          leading:Radio(
                            groupValue: rescuePostLevel,
                            value:'B',
                            onChanged:(String? value){
                              setState((){
                                rescuePostLevel = value;
                              });
                            },
                          ),
                          onTap:(){
                            setState((){
                              rescuePostLevel = 'B';
                            });
                          },
                        ),
                        ListTile(
                          title:Text('委員長レベル'),
                          leading:Radio(
                            groupValue: rescuePostLevel,
                            value: 'C',
                            onChanged:(String? value){
                              setState((){
                                rescuePostLevel = value;
                              });
                            },
                          ),
                          onTap:(){
                            setState((){
                              rescuePostLevel = 'C';
                            });
                          },
                        ),
                        Text('緊急度',style:Theme.of(context).textTheme.titleLarge),
                        Text('必要な対応速度を選択してください'),
                        ListTile(
                          title:Text('10分以内'),
                          leading:Radio(
                            groupValue: rescuePostLeve2,
                            value: 'A',
                            onChanged:(String? value){
                              setState((){
                                rescuePostLeve2 = value;
                              });
                            },
                          ),
                          onTap:(){
                            setState((){
                              rescuePostLeve2 = 'A';
                            });
                          },
                        ),
                        ListTile(
                          title:Text('30分以内'),
                          leading:Radio(
                            groupValue: rescuePostLeve2,
                            value: 'B',
                            onChanged:(String? value){
                              setState((){
                                rescuePostLeve2 = value;
                              });
                            },
                          ),
                          onTap:(){
                            setState((){
                              rescuePostLeve2 = 'B';
                            });
                          },
                        ),
                        ListTile(
                          title:Text('それ以上'),
                          leading:Radio(
                            groupValue: rescuePostLeve2,
                            value: 'C',
                            onChanged:(String? value){
                              setState((){
                                rescuePostLeve2 = value;
                              });
                            },
                          ),
                          onTap:(){
                            setState((){
                              rescuePostLeve2 = 'C';
                            });
                          },
                        ),
                      ]
                  )

                ),
                //3ページ目のcontent
                SingleChildScrollView(
                  child:ListBody(
                    children:<Widget>[
                      Text('本当に送信しますか'),
                    ]
                  )
                ),
                SingleChildScrollView(
                   child:(
                    Text('送信しました')
                  )
                ),
              ]
            )
            ),
            //条件分岐でactionを記述
            actions:<Widget>[
              if(rescuePageIndex == 0)...[
                TextButton(
                  child:
                    Text('キャンセル'),
                    onPressed: () {
                      Navigator.of(context).pop();
                    }
                ),
                TextButton(
                  child: 
                    Text('次へ'),
                    onPressed: () {
                      setState((){
                        if (manualChecked==true && poepleChecked==true){
                          rescuePageIndex = 1;
                          confirmError1 = false;
                        }
                        else{
                          confirmError1 = true;
                        }
                      });
                    }
                )
              ]
              else if (rescuePageIndex == 1)...[
                TextButton(
                  child:
                    const Text('戻る'),
                      onPressed: () {
                        setState((){
                        rescuePageIndex = 0;
                        });
                      }
                ),
                TextButton(
                  child: 
                    const Text('送信'),
                    onPressed: () {
                      setState((){
                        if (taskChecked==true && locateChecked==true){
                          rescuePageIndex = 2;
                          confirmError2 = false;
                        }
                        else{
                          confirmError2 = true;
                        }
                      });
                    }
                )
              ]
              else if (rescuePageIndex == 2)...[
                TextButton(
                  child:(
                    Text('キャンセル')
                  ),
                  onPressed:(){
                    Navigator.of(context).pop();
                  }
                ),
                TextButton(
                  child:(
                    Text('送信')
                  ),
                  onPressed:(){
                    // setState((){
                    //   api.postRescue();
                    //   _rescuePageIndex = 3;//データを送信する機能はここに入れる
                    // });
                  }
                )
              ]
              else if (rescuePageIndex == 3)
                TextButton(
                  child:(
                    Text('閉じる')
                  ),
                  onPressed:(){
                    Navigator.of(context).pop();
                    }
                  )
            ]
          );
        }
      );
    },
  );
}
//以下削除予定
          // IndexedStack(
          //   index: _rescuePageIndex,
          //   children:[
          //     //1ページ目
          //     AlertDialog(
          //       title:  Text('レスキュー要請1/2'),
          //       content: SingleChildScrollView(
          //         child: ListBody(
          //           children: <Widget>[
          //             Text('確認', style:Theme.of(context).textTheme.titleLarge),
          //             Text('本部に応援を要請する前に、以下のことを行ってください'),
          //             Material(
          //               color:Colors.transparent,
          //               child:InkWell(
          //                 splashColor:Colors.transparent,
          //                 highlightColor:Colors.transparent,
          //                 onTap:(){
          //                       setState((){
          //                         _manualChecked = !_manualChecked;
          //                       });
          //                     },
          //                 child:Row(
          //                   children:[
          //                     Checkbox(
          //                       value: _manualChecked,
          //                       onChanged: (bool? value) {
          //                         setState(() { // 🔹 `setState()` で UI を更新
          //                           _manualChecked = value!;
          //                         });
          //                       },
          //                     ),
          //                     Text('マニュアルを確認した')
          //                   ],
          //                     ),
          //                 ),
          //             ),
          //             Material(
          //               color:Colors.transparent,
          //               child:InkWell(
          //                 splashColor:Colors.transparent,
          //                 highlightColor:Colors.transparent,
          //                 onTap:(){
          //                       setState((){
          //                         _poepleChecked= !_poepleChecked;
          //                       });
          //                     },
          //               child:Row(
          //                 children:[
          //                   Checkbox(
          //                   value: _poepleChecked,
          //                   onChanged: (bool? value) {
          //                     setState(() { // 🔹 `setState()` で UI を更新
          //                       _poepleChecked = value!;
          //                     });
          //                   },
          //                   ),
          //                   Text('近くの先輩に相談をした')

          //                 ]
          //               )
          //               ),
          //             ),
          //             if (_confirmError1)
          //               Text('確認をしてください',style:TextStyle(color:Colors.red))
          //           ],
          //         ),
          //       ),
          //       actions: <Widget>[
          //         TextButton(
          //           child:
          //             const Text('キャンセル'),
          //               onPressed: () {
          //                 Navigator.of(context).pop();
          //               }
          //         ),
          //         TextButton(
          //           child: 
          //             const Text('次へ'),
          //               onPressed: () {
          //                 setState((){
          //                   if (_manualChecked==true && _poepleChecked==true){
          //                     _rescuePageIndex = 1;
          //                     _confirmError1 = false;
          //                   }
          //                   else{
          //                     _confirmError1 = true;
          //                   }
          //                 });
          //               }
          //         )
          //       ],
          //     ),
          //     //2ページ目
          //     AlertDialog(
          //       title: Text('レスキュー要請2/2'),
          //       content: SingleChildScrollView(
          //         child:ListBody(
          //             children:<Widget>[
          //               Text('タスク', style:Theme.of(context).textTheme.titleLarge),
          //               Padding(
          //                 padding: EdgeInsets.all(16.0),
          //                 child:
          //                 Text('問題の発生したタスクを選択してください'),
          //                 ),
          //               Autocomplete<String>(
          //                 optionsBuilder: (TextEditingValue textEditingValue) {
          //                   List<String> tasklist = ['駐車場誘導', '受付（講義棟前）', '受付（プール前）', 'その他'];

          //                   // 🔍 ユーザーの入力に基づいてリストをフィルタリング
          //                   return tasklist.where((String item) {
          //                     return item.toLowerCase().contains(textEditingValue.text.toLowerCase());
          //                   }).toList();
          //                 },
          //                 onSelected: (String selecttask) {
          //                   _taskChecked = !_taskChecked;
          //                 },
          //                 fieldViewBuilder: (context, textEditingController, focusNode, onFieldSubmitted) {
          //                   return TextField(
          //                     controller: textEditingController,
          //                     focusNode: focusNode,
          //                     decoration: InputDecoration(
          //                       labelText: "タスク",
          //                       prefixIcon: Icon(Icons.search), // 🔍 検索アイコンを追加
          //                       border: OutlineInputBorder(),
          //                     ),
          //                   );
          //                 },
          //               ), 
          //               if(_taskChecked == false && _confirmError2 == true)
          //                 Text('必要事項を入力してください',style:TextStyle(color:Colors.red)),
          //               SizedBox(height: 20),
          //               Text('場所', style:Theme.of(context).textTheme.titleLarge),
          //               Padding(
          //                 padding: EdgeInsets.all(16.0),
          //                 child:
          //                 Text('現在の場所を教えてください'),
          //                 ),
          //               Autocomplete<String>(
          //                 optionsBuilder: (TextEditingValue textEditingValue) {
          //                   List<String> locatelist= ['体育館','プール前','講義棟前','駐車場入口','電気棟駐車場'];
          //                   return locatelist.where((String item){
          //                     return item.toLowerCase().contains(textEditingValue.text.toLowerCase());
          //                   }).toList();
          //                 },
          //                 onSelected: (String selectlocate) {
          //                   _locateChecked = !_locateChecked;
          //                 },
          //                 fieldViewBuilder: (context, textEditingController, focusNode, onFieldSubmitted) {
          //                   return TextField(
          //                     controller: textEditingController,
          //                     focusNode: focusNode,
          //                     decoration: InputDecoration(
          //                       labelText: "場所",
          //                       prefixIcon: Icon(Icons.search), // 🔍 検索アイコンを追加
          //                       border: OutlineInputBorder(),
          //                     ),
          //                   );
          //                 }
          //               ),
          //               if(_locateChecked == false && _confirmError2 == true)
          //                 Text('必要事項を入力してください',style:TextStyle(color:Colors.red)),
          //               Text('内容', style:Theme.of(context).textTheme.titleLarge),
                        
          //             ]
          //         )
          //       ),
          //       actions: <Widget> [
          //         TextButton(
          //           child:
          //             const Text('戻る'),
          //               onPressed: () {
          //                 setState((){
          //                 _rescuePageIndex = 1;
          //                 });
          //               }
          //         ),
          //         TextButton(
          //           child: 
          //             const Text('送信'),
          //             onPressed: () {
          //               setState((){
          //                 if (_taskChecked==true && _locateChecked==true ){
          //                   _rescuePageIndex = 2;
          //                   _confirmError2 = false;
          //                 }
          //                 else{
          //                   _confirmError2 = true;
          //                 }
          //               });
          //             }
          //         )
          //       ]
          //     ),
          //     //3ページ目
          //     AlertDialog(
          //       title:Text('最終確認'),
          //       content:SingleChildScrollView(
          //         child:ListBody(
          //           children:<Widget>[
          //             Text('本当に送信しますか'),
          //           ]
          //         )
          //       ),
          //       actions:<Widget>[
          //         TextButton(
          //           child:(
          //             Text('キャンセル')
          //           ),
          //           onPressed:(){
          //             Navigator.of(context).pop();
          //           }
          //         ),
          //         TextButton(
          //           child:(
          //             Text('送信')
          //           ),
          //           onPressed:(){
          //             setState((){
          //               _rescuePageIndex = 3;//データを送信する機能はここに入れる
          //             });
          //           }
          //         )
          //       ]
          //     ),
          //     //4ページ目
          //     AlertDialog(
          //       content:SingleChildScrollView(
          //         child:(
          //           Text('送信しました')
          //         )
          //       ),
          //       actions:<Widget>[
          //         TextButton(
          //           child:(
          //             Text('閉じる')
          //           ),
          //           onPressed:(){
          //             Navigator.of(context).pop();
          //           }
          //         )
          //       ]
          //     )
          //   ]
          // );