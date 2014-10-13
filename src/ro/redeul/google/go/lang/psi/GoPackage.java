package ro.redeul.google.go.lang.psi;

import java.util.Collection;

public interface GoPackage extends GoPsiElement {

    public String getImportPath();

    public String getName();

    Collection<GoFile> getPackageFiles();
}
