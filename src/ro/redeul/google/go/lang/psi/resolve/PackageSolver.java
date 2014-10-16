package ro.redeul.google.go.lang.psi.resolve;

import com.intellij.psi.PsiElement;
import ro.redeul.google.go.lang.psi.GoPackage;
import ro.redeul.google.go.lang.psi.GoPackageReference;
import ro.redeul.google.go.lang.psi.resolve.refs.PackageReference;
import ro.redeul.google.go.lang.psi.toplevel.GoImportDeclaration;

public class PackageSolver extends GoReferenceSolver<PackageReference, PackageSolver> {


    public PackageSolver(PackageReference reference) {
        super(reference);
    }

    @Override
    public void visitImportDeclaration(GoImportDeclaration declaration) {
        if (isReferenceTo(declaration))
            addTarget(declaration);
    }

    boolean isReferenceTo(GoImportDeclaration importDeclaration) {

        GoPackageReference packageReference = importDeclaration.getPackageReference();

        String packageName = null;
        if ( packageReference != null && !(packageReference.isBlank() || packageReference.isLocal()) )
            packageName = packageReference.getString();
        else {
            GoPackage goPackage = importDeclaration.getPackage();
            packageName = goPackage != null ? goPackage.getName() : null;
        }

        return packageName != null && packageName.equals(getReference().getCanonicalText());
    }

    @Override
    public PsiElement resolve(PackageReference reference) {
        return null;
    }

    @Override
    public Object[] getVariants() {
        return new Object[0];
    }
}
