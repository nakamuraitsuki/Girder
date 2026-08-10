package infrastructure

// Package infrastructure では、Girderの裏で動くDriver関連の実装を提供する。
// 基本方針として、このプロジェクトは裏にあるLibvirtやOVNの挙動をオーけストレートする目的なので、
// 過剰な抽象化やDTOの乱立は避ける方針である。よってこのInfrastructureとAPIの間には、Infrastructureを用いない。
// 依存性逆転の意味がないので。