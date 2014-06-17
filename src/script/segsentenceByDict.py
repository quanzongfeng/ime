#/***************************************************************************
# * 
# * Copyright (c) 2014 Baidu.com, Inc. All Rights Reserved
# * 
# **************************************************************************/
# 
# 
# 
#/**
# * @file segsentenceByDict.py
# * @author quanzongfeng(com@baidu.com)
# * @date 2014/05/12 11:50:29
# * @brief 
# *  
# **/

import sys

SysDict={}
def loadterms(term):
    f = open(term)
    while 1:
        a = f.readline(-1)
        if a == '':
            break
        b = a.rstrip().split()
        if len(b) < 3:
            continue
        hz = b[0]
        freq = b[1]
        SysDict.setdefault(hz, freq)

        if freq > SysDict[hz]:
            SysDict[hz] = freq

    f.close()

def isHz(hz, length):
    ln = len(hz)
    if length <= 0: 
        return 0
    if ln == 0:
        return 0
    if ord(hz[0]) <= 128:
        return 1
    if ord(hz[0])> 128 && length == 1:
        return -1
    return 2

def getIndexList(sent):
    tn = len(sent)
    i = 0
    while (i< tn):



def segForce(sent, dict):
    ln = len(sent)


            























#/* vim: set expandtab ts=4 sw=4 sts=4 tw=100: */
