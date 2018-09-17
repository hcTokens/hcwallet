#include <stdio.h>
#include <stdio.h>
#include <windows.h>
#include "omniproxy.h"


typedef char* (WINAPI *FunJsonCmdReq)(char *);
typedef int (WINAPI *FunOmniStart)(char *);
typedef int (WINAPI *FunSetCallback)(unsigned int,void *);

FunOmniStart    funOmniStart = NULL; //
FunSetCallback  funSetCallback=NULL;
FunJsonCmdReq   funJsonCmdReq = NULL; //
#define INDEX_CALLBACK_GoJsonCmdReq 1


//extern char* GJsonCmdReq(char *pcReq);

void CLoadLibAndInit()
{
	printf("in LoadDllStart\n");

	HINSTANCE hDllInst = LoadLibrary("omnicored.DLL");
    if(!hDllInst)
    {
        //FreeLibrary(hDllInst);
        return;
    }

    funOmniStart = (FunOmniStart)GetProcAddress(hDllInst,"OmniStart");
    funJsonCmdReq= (FunJsonCmdReq)GetProcAddress(hDllInst,"JsonCmdReq");
    funSetCallback= (FunSetCallback)GetProcAddress(hDllInst,"SetCallback");
    //以后如果需要,在dll里面加SetFunc(index,Func)，序号,函数,类似接口
    printf("funJsonCmdReq=%d",funJsonCmdReq);


    if(funSetCallback) //先setcallback后 omnistart
       funSetCallback(INDEX_CALLBACK_GoJsonCmdReq,JsonCmdReqOmToHc);


    // myFunOmniStart 在DLL中声明的函数名
    //if(funOmniStart)
    //    funOmniStart(pcArgs);

    /*
	dll := syscall.MustLoadDLL("omnicored.dll")
	procGreet := dll.MustFindProc("OmniStart")
	procGreet.Call()
	*/

    return;
}

int COmniStart(char *pcArgs)
{
    if(funOmniStart==NULL)
        return -1;
    return funOmniStart(pcArgs);
}

char* CJsonCmdReq(char *pcReq)
{
    if(funJsonCmdReq==NULL)
        return NULL;
    return funJsonCmdReq(pcReq);
};

int CSetCallback(int iIndex,void* pCallback)
{
    if(funSetCallback==NULL) return -1;
    if(pCallback==NULL) return -1;
    return funSetCallback(iIndex,pCallback);
};

