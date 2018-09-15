#include <GoInclude.h>

#include <windows.h>
#include <ImportDll.h>

typedef UINT(CALLBACK* lpCallback)(int, void*);
lpCallback lpfnDllApi = NULL;

int DllCallBack(int nType)
{
	return GoCallback(nType);
}

char* DllCallBackChar(int nType, char* str)
{
	return GoCallbackChar(nType, str);
}

char* DllCallBackCharEx(int nType, char* str, int len)
{
	return GoCallbackCharEx(nType, str, len);
}

int initDll()
{
	HINSTANCE hDLL = LoadLibrary("bitcoind_dll.dll");
	if (hDLL != NULL)
	{
		lpfnDllApi = (lpCallback)GetProcAddress(hDLL, "BtcStart");
		if (!lpfnDllApi)
		{
			// handle the error
			FreeLibrary(hDLL);
			return -2;
		}
		else
		{
			// call the function
			lpfnDllApi(0, (void*)DllCallBack);
			lpfnDllApi(1, (void*)DllCallBackChar);
			lpfnDllApi(2, (void*)DllCallBackCharEx);
			return 0;
		}
	}
	return -1;
}
int CallCpp(int type, void* param)
{
	switch (type)
	{
	case 0:
		return initDll();
		break;
	case 1:
		if (lpfnDllApi)
			lpfnDllApi(3, (void*)param);
		break;
	case -1:
	{
		/*
			char temp[] = {"[{\"fromaddress\":\"mnywEASo8FssNNhXm44hDVeZwMgZbH7suy\",\"payload\":\"\\u0000\\u0000\\u00002\\u0002\\u0000\\u0001\\u0000\\u0000\\u0000\\u0000Companies\\u0000Bitcoin Mining\\u0000Quantum Miner\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u000fB@\"}]"};
			DllCallBackChar(12, temp);

			char temp[] = {"fromaddress===TsSuLQReXq4TDf6mDgiN3j6rLHQpF73hQQD;;;payload===omni\0\0\02\x2\0\x1\0\0\0\0Companies\0Bitcoin Mining\0Quantum Miner\0\0\0\0\0\0\0\0\xfB@"};
			GoCallbackCharEx(12, temp, 126);
			*/
		char temp[] = { "TsbXqRS4p7M8EzJ2nXur3S9sA3zkqziXZv1" };
		GoCallbackCharEx(13, temp, 35);
	}
	break;
	case 4:
		if (lpfnDllApi)
			lpfnDllApi(4, (void*)param);
		break;
	default:
		break;
	}
	return -1;
	//return(callback2("abc"));
};

