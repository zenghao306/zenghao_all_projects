#ifdef WIN32
#if defined(EXPORT_DLL)
#    define VAR __declspec(dllexport)
#elif defined(IMPORT_DLL)
#    define VAR __declspec(dllimport)
#endif
#else
#    define VAR extern
#endif

extern int hdt_encode_v0(const char* src, char** encoded_dst);
extern int hdt_decode_v0(const char* encoded_src, char** decoded_dst);
extern void hdt_release(void* ptr);
