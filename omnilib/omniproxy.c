#include <stdio.h>
#include <stdio.h>
#include <windows.h>
#include "omniproxy.h"


typedef const char* (WINAPI *FunJsonCmdReq)(char *);
typedef int (WINAPI *FunOmniStart)(char *);
typedef int (WINAPI *FunSetCallback)(unsigned int,void *);

FunOmniStart    funOmniStart = NULL; //
FunSetCallback  funSetCallback=NULL;
FunJsonCmdReq   funJsonCmdReq = NULL; //
#define INDEX_CALLBACK_GoJsonCmdReq 1

void CLoadLibAndInit()
{
	printf("in LoadDllStart\n");

	HINSTANCE hDllInst = LoadLibrary("omnicoind.dll");
    if(!hDllInst)
    {
        //FreeLibrary(hDllInst);
        return;
    }

    funOmniStart = (FunOmniStart)GetProcAddress(hDllInst,"OmniStart");
    funJsonCmdReq= (FunJsonCmdReq)GetProcAddress(hDllInst,"JsonCmdReq");
    funSetCallback= (FunSetCallback)GetProcAddress(hDllInst,"SetCallback");

    printf("funJsonCmdReq=%d",funJsonCmdReq);
    return;
}

int COmniStart(char *pcArgs)
{
    if(funOmniStart==NULL)
        return -1;
    return funOmniStart(pcArgs);
}

const char* CJsonCmdReq(char *pcReq)
{
    if(funJsonCmdReq==NULL)
        return NULL;
    const char* ret = funJsonCmdReq(pcReq);
    printf("88888888888888888888888888888\n");
    printf(ret);
    return ret;
};

int CSetCallback(int iIndex,void* pCallback)
{
    if(funSetCallback==NULL) return -1;
    if(pCallback==NULL) return -1;
    return funSetCallback(iIndex,pCallback);
};

