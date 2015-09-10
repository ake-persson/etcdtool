%define name %NAME%
%define version %VERSION%
%define release %RELEASE%
%define buildroot %{_topdir}/BUILDROOT
%define sources %{_topdir}/SOURCES

BuildRoot: %{buildroot}
Source: %SOURCE%
Summary: %{name}
Name: %{name}
Version: %{version}
Release: %{release}
License: Apache License, Version 2.0
Group: System
AutoReqProv: no

%description
%{name}

%prep
mkdir -p %{buildroot}/usr/bin
cp %{sources}/bin/* %{buildroot}/usr/bin

%files
%defattr(-,root,root)
/usr/bin/*
