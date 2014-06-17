#/***************************************************************************
# * 
# * Copyright (c) 2014 Baidu.com, Inc. All Rights Reserved
# * 
# **************************************************************************/
# 
# 
# 
#/**
# * @file segSentence.py
# * @author quanzongfeng(com@baidu.com)
# * @date 2014/04/24 16:28:40
# * @brief 
# *  
# **/

import sys
import urllib
import json
import socket
import time

def segsentence(line):
    url = "http://nlpc.baidu.com/?method=wordseg&encoding=gbk&query="
    url += line
#    print url

    hz_basic_list= []
    hz_seg_list= []
    t = 0
    tt = ''
    try :
        t = urllib.urlopen(url)
        bf = t.read()

        tt = json.loads(bf, 'gbk')
    except:
        if t!= 0:
            t.close()
        return [],[]

    ind = 0
    v = tt['BasicWordResult']
    for p in v:
        hz = p['buffer'].encode('gbk')
        of = p['offset']
        if int(of) >= len(hz_basic_list):
            hz_basic_list = hz_basic_list + [""]*(int(of)+1)
            hz_basic_list[int(of)] = hz
        else:
            hz_basic_list[int(of)] = hz


    vs = tt['SegmentResult']
    ind = 0
    for p in vs:
        hz = p['buffer'].encode('gbk')
        of = p['offset']
        if int(of) >= len(hz_seg_list):
            hz_seg_list = hz_seg_list + [""]*(int(of)+1)
            hz_seg_list[int(of)] = hz
        else:
            hz_seg_list[int(of)] = hz

    t.close()

    #print hz_basic_list
    #print hz_seg_list
    basic_list = [t for t in hz_basic_list if t != '']
    seg_list = [t for t in hz_seg_list if t != '']
    return basic_list, seg_list

def isHz(line):
    if len(line) == 0:
        return 0, ''
    if ord(line[0]) > 128:
        if len(line) == 1:  #error
            return 0, ''
        return 2, line[0:2]
    else:
        return 1, line[0:1]

def read_terms(file):
    f = open(file)
    word_dict = {}
    while 1:
        a = f.readline()
        if a == '':
            break
        b =a.rstrip().split()
        hz = b[0]

        word_dict.setdefault(hz, 1)

    return word_dict

def filter_list(ll, dt):
    for t in ll:
        dt.setdefault(t, 0)
        if dt[t] == 0:
            return 0
    return 1


def segment(sent, dt):
    bas, sl = segsentence(sent)
    reseg = []
    rebas = []

    if len(sl) != 1 and filter_list(sl, dt) == 1:
        reseg = sl
        for t in reseg :
            if len(t) < 6:
                return reseg, []
    if len(rebas) != 1 and  filter_list(bas, dt) == 1:
        rebas = bas
    return reseg, rebas


def read_lines(term):
    dt = read_terms(term)
    socket.setdefaulttimeout(5)
    i = 0
    while 1:
        f = sys.stdin.readline()
        if f == '':
            break

        tokens = f.rstrip().split()
        sent = tokens[1]

        seg, bas = segment(sent, dt)
        flag = 0
        if seg != []:
            flag |= 4 
        if bas != []:
            flag |= 1

        result = "\t".join(tokens)
        num = 0
        if flag > 4:
            num = 2
        elif flag == 4:
            num = 1
        elif flag == 1:
            num = 1
        else:
            num = 0

        #print seg, bas, flag, num
        result += "\t"+str(num)
        if flag & 4 != 0 :
            result += "\t" + str(len(seg))+"\t"+"\t".join(seg)
        if flag & 1 != 0:
            result += "\t" + str(len(bas))+"\t"+"\t".join(bas)

        print result


if __name__=='__main__':
    if len(sys.argv) == 1:
        print "need term dict path"
        sys.exit(1)
    read_lines(sys.argv[1])




    
        










        












#/* vim: set expandtab ts=4 sw=4 sts=4 tw=100: */
