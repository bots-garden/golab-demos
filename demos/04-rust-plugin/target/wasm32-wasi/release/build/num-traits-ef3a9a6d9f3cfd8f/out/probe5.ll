; ModuleID = 'probe5.34344da6fcc24e5d-cgu.0'
source_filename = "probe5.34344da6fcc24e5d-cgu.0"
target datalayout = "e-m:e-p:32:32-p10:8:8-p20:8:8-i64:64-n32:64-S128-ni:1:10:20"
target triple = "wasm32-unknown-wasi"

@alloc_f421969b09b54afe21712b0d66a11b5a = private unnamed_addr constant <{ [77 x i8] }> <{ [77 x i8] c"/rustc/5680fa18feaa87f3ff04063800aec256c3d4b4be/library/core/src/ops/arith.rs" }>, align 1
@alloc_657e755a5a223ff4bac9087427e3017e = private unnamed_addr constant <{ ptr, [12 x i8] }> <{ ptr @alloc_f421969b09b54afe21712b0d66a11b5a, [12 x i8] c"M\00\00\00\F8\02\00\00\01\00\00\00" }>, align 4
@str.0 = internal constant [28 x i8] c"attempt to add with overflow"
@alloc_2e38410fced2c310c68bdf2d45d0c3bd = private unnamed_addr constant <{ [4 x i8] }> <{ [4 x i8] c"\02\00\00\00" }>, align 4

; <i32 as core::ops::arith::AddAssign<&i32>>::add_assign
; Function Attrs: inlinehint nounwind
define internal void @"_ZN66_$LT$i32$u20$as$u20$core..ops..arith..AddAssign$LT$$RF$i32$GT$$GT$10add_assign17h038489792c8f6a40E"(ptr align 4 %self, ptr align 4 %other) unnamed_addr #0 {
start:
  %other1 = load i32, ptr %other, align 4, !noundef !0
  %0 = load i32, ptr %self, align 4, !noundef !0
  %1 = call { i32, i1 } @llvm.sadd.with.overflow.i32(i32 %0, i32 %other1)
  %_4.0 = extractvalue { i32, i1 } %1, 0
  %_4.1 = extractvalue { i32, i1 } %1, 1
  %2 = call i1 @llvm.expect.i1(i1 %_4.1, i1 false)
  br i1 %2, label %panic, label %bb1

bb1:                                              ; preds = %start
  store i32 %_4.0, ptr %self, align 4
  ret void

panic:                                            ; preds = %start
; call core::panicking::panic
  call void @_ZN4core9panicking5panic17h2f7f1e15b374dfd2E(ptr align 1 @str.0, i32 28, ptr align 4 @alloc_657e755a5a223ff4bac9087427e3017e) #5
  unreachable
}

; probe5::probe
; Function Attrs: nounwind
define hidden void @_ZN6probe55probe17h60e6f74554da1bf3E() unnamed_addr #1 {
start:
  %x = alloca i32, align 4
  store i32 1, ptr %x, align 4
; call <i32 as core::ops::arith::AddAssign<&i32>>::add_assign
  call void @"_ZN66_$LT$i32$u20$as$u20$core..ops..arith..AddAssign$LT$$RF$i32$GT$$GT$10add_assign17h038489792c8f6a40E"(ptr align 4 %x, ptr align 4 @alloc_2e38410fced2c310c68bdf2d45d0c3bd) #6
  ret void
}

; Function Attrs: nocallback nofree nosync nounwind speculatable willreturn memory(none)
declare hidden { i32, i1 } @llvm.sadd.with.overflow.i32(i32, i32) #2

; Function Attrs: nocallback nofree nosync nounwind willreturn memory(none)
declare hidden i1 @llvm.expect.i1(i1, i1) #3

; core::panicking::panic
; Function Attrs: cold noinline noreturn nounwind
declare dso_local void @_ZN4core9panicking5panic17h2f7f1e15b374dfd2E(ptr align 1, i32, ptr align 4) unnamed_addr #4

attributes #0 = { inlinehint nounwind "target-cpu"="generic" }
attributes #1 = { nounwind "target-cpu"="generic" }
attributes #2 = { nocallback nofree nosync nounwind speculatable willreturn memory(none) }
attributes #3 = { nocallback nofree nosync nounwind willreturn memory(none) }
attributes #4 = { cold noinline noreturn nounwind "target-cpu"="generic" }
attributes #5 = { noreturn nounwind }
attributes #6 = { nounwind }

!0 = !{}
