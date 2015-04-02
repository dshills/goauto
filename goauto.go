// Copyright 2015 Davin Hills. All rights reserved.
// MIT license. License details can be found in the LICENSE file.

// Package goauto implements a set of tools for building workflow automation tools.
// These tools can be as simple as running a compiler when a source file changes to complex chains of tasks doing almost any action required within a development environment
// See README.md for more details on usage
package goauto

// Verbose is a global var that will print a lot of debug info during processing
// This is handy for debugging. By default it is off
/*

TODO
More built in tasks
Write more tests. Can always use more tests
Call your mother

- War Doctor
                                         ......
                                        .::,=?~,...
                                       .,+?$?,?=,=.
                                  .  ..=::.~II.~~.:.
                                   .:?8O.,~,,?+??,I~,,.
                                 .,~7+I=,.~?I7I~?II=?:.
                               ....,,,+I:Z$?,+D7,Z,+:~,,.
                               ,::$~?+,?,I,?+=I+8,8~~..,.
                             .,:.7+$:7$~+~?$D$I7+O==,~+,:.
                            .I$=~+OZI7I7777ZZ$$ZZZ$7+~I+,:,
                           ..+.~,=?+=+7$7I$I$77$ZZZ$==:=..,.
                            ,..+.+?==I7$7ZO8OZZZZZO$I+=7=:,.
                            .,~Z~.?==++7$ZOO8OOOZ$Z7II++~~:..
                           ..:?.::+=+7$$Z7OZZZ7Z7ZOO$+=7,.~.
                           .,.7=:+=+?7$I?$7II7I?IZZZZ7?I7:=,
                            .,+=++=:?III$I~~=+??+?$Z$7?+7?..
                            ,:~:?~~,~???O7$O+7Z77ZDOI==:77.:.
                           .:?===:~~~I7?7D88?7OZ88N$ZII,+~++
                           =I~:=+,,,~+?+7$87??Z+O8O7$7=,:,??=
                          .Z++==:,~~+++I7Z7I?IOOIZ$$I+~.,,=I7
                           7~I?++,.:+?$$I7?=~IZZ:ZZ$Z7+~,~I7I
                           ,~Z$?:,,:?+I?IIZ==:==,Z$$I?+,.I7O.
                            7$ZI=.,~===+?7$7IZZZ$$$I~I+.,O?.
                            .~=?=.:====+$?7$I++??I?7~I=.==~
                             .~?7.::+I+=I+?II?+?++~.~I=,??
                           .88:,~.~:==~~II??I+7I=?=~.+:.
                          .DNN.Z=.,,++:=~?I7I7I7Z7+~,~,,
                        .78NNN~O7,.~:==~+?~=+I7+???+.,~=
                       =DODNNN?+OO~~,~=+?I+?+???7$7I+,,N8.
                      O8ODNNND?=I=N$I,=~+==???++=7+:,,?IDD=
                     .8D8DNNND877O$ZOI,::=:+=I??+=:.=OZ+DDDZ
                     I8M8DNNNDNZDD=N78?~~~=~++~~=::+8$OO7DDD$
                    ?DD8DDNNNDNDD7DDN8?Z+I=~II~~+~+N8DO8.DZD8:
                =7$ODD88ONDNNDDN+88$DO8778O?II=+7+?DNZ8$=DDDZZ.
          . ~?7ZO88ND8DD8DDNNDDNN$Z$$ONND$D8$I$$8INDNN$DI8DD88$..
       ?ZZZZDD88D8NNMDD88DDNNDDDNNOOZ+D$8D8D$8ODOONM8ONZ$NN8DZ8:.=?.
    $DZ88O8DDNDDDDNN8NDDDDDDNDDDNNDI$N=8NZNOONNDDMNN8ONZINNNNNDO?:.:.
  IDD8$D8D8DDNNDNNNNZ8DDDDDDNDDNDNN88ODZOMN8DDNNNMMNOONDINNNNDND8.$....
  DDDDDDD8DDNNNNDNNNONDNNDDDDDDNDNMM8OZ8$Z8O$N8NNNN8NMDI?NMNNNDDN8DZ.. ..
  8NNDD88DDDNDNDDDNNONDDDDDDDNNNNNNMMN8IIODNO$ZI888+MND+NNNNNZNDZOO?~,.=.
 ONDNDDDN8NDNNDNDD8N8NNDNNDNNNNDNNNMMMMD$88OO$O+NOND7D,OMMMNNDODD87Z$,~:O.
  8DDD8D8O88DDDDDDO88NNNNNDNDNNNNNNNMNNN8+8OZ$D~N8DND$8NNMNMND$ND$8Z?7:,:+,
 D8DN8D88888DNNDDD8OODNNDDDNDNONNDNNMMNNNZD+7?ODDDN88IINNMNNNN=8N87,?$,~:=+.
 DNDDON$ZDZDDDNNND8OONNN8DMMMMNDDNNNMMNMMO8?+?8D7DDDOODNNNNMMND8N8=~.=O:7+II.
 DDMDD88ODO88NNND8OOOOZ8NNMMNNDNNNNNNMMNNM$+I+I88?NND7:DNNMMMNNDDNO+,~~$,~$D$
Z8NDN8NDZNODDNNNDDOOOZZZDN88NNDDNDNNNNMMNDND7NZ8?$?DZDOIDMDNDNNINN$7=,,,~=ZD8
Z88DNND8ZDZDDNNNDDDOO8ZZODDDDODNNNNNNNNMNNNMI8ID$78$D8DZ8NM$NNN:NND7::~=.~$ND
O88DNNNND8O88NNNDD888O8D8ON88DNNNDNNNNNNMMNDN8O$ODZ:OD8?7MMMNNMZNNNZ.:::.I~DD~
 8DDDNNN8ZZODDNNNDD8OZ8OD8D8NDDDNDNDNNNNNNMNDDDZZD88$DO7?+8DDMNNNNNDI=.~:Z~DDZ
Z88DDDNM88Z$DNNNND8$8ODO88$NNNNNDNDDNNNNNMMMNNND8=78O8DO=$NNNMNNONNDO?:::Z?8D8
  8DDNNNN8ZO8DNNNN88OZ8DODMN8DNDNDNDNNNNNNNNNNNDD?OZ$$$Z8$$8NMMN.MND8I,+~8Z=$D
  NDDDNNN8Z$ODNNNND88OIOODN8ZNNDDDNDDNNNDDNNMMNNNDMMNDMN$N8NOMMNONNNDI:+=DD+DD
  ONDNNNNNOZO8DNNND8OZZZO7$8NN8NN8NNMNDNNNDNMMMMN88MNNN88DDOOMNNNNNNNZ+$=O87D+.
  DDDNNNNNO$O8DNNND8OOOZZOON$8NNN8NNNNN8NNDNMMMMN$DNNMMO88D$OMNMN,NNNO+$+Z8$78:
   8NNNNMMOZZODNNN8D8OOZOOODNNNO8NONNNNNDNN8NMMNNNIDNM8D8DN?MMNNN7NND887Z788+D?
   88DNNMMDZ$Z8DNNDDOZZZ$ZODNMNDO8$ODNNNNNNN8NMMMMNNN8O8DDOMDMMMN8NNNNN8O?88?D+
   D8NNNNNNOZZ8NNNDD8OZZZZZDNMNDOODDN$NNNNNNNDDMMDMM88DDNNNMNNMNNNNNNNDOD$$D7DZ
    DDNNNNMO$Z8NNNDD8ZZZZZODDNNDOZDDNNMM$$8DNNDONMMMDDDN8MMMMMMMND=NND+$8ZDD8DZ
    8NNNNNMO$ZONNND88ZZ$$$8DNNNDOZ8DNDNMMDD8OOD8MMN8DNNDMNMMMNNNND$MND:=OO8D$D8
     DDNNNNOZZ8DNNDD8ZZ$ZZZDNMNDZZ8DNNNMNNDDDOODIOOO8NNNMMMNNNMMNNNND8==7ZD87D8
*/
