package ro.redeul.google.go.lang.psi.resolve.references;

import com.intellij.patterns.ElementPattern;
import com.intellij.psi.PsiElement;
import com.intellij.psi.impl.source.resolve.ResolveCache;
import org.jetbrains.annotations.NotNull;
import ro.redeul.google.go.lang.psi.expressions.literals.GoLiteralIdentifier;
import ro.redeul.google.go.lang.psi.expressions.primary.GoLiteralExpression;
import ro.redeul.google.go.lang.psi.processors.ResolveStates;
import ro.redeul.google.go.lang.psi.resolve.GoResolveResult;
import ro.redeul.google.go.lang.psi.resolve.PackageSolver;
import ro.redeul.google.go.lang.psi.resolve.ResolveCacheResolvers;
import ro.redeul.google.go.lang.psi.utils.GoPsiScopesUtil;

import static com.intellij.patterns.PlatformPatterns.psiElement;

public class PackageReference extends Reference<GoLiteralIdentifier, GoLiteralIdentifier, PackageSolver, PackageReference> {

    public static ElementPattern<GoLiteralIdentifier> MATCHER =
            psiElement(GoLiteralIdentifier.class)
                    .withParent(
                            psiElement(GoLiteralExpression.class)
                    );

    private static final ResolveCache.AbstractResolver<PackageReference, GoResolveResult> RESOLVER =
            ResolveCacheResolvers.<PackageReference, PackageSolver>makeDefault();

    public PackageReference(@NotNull GoLiteralIdentifier element) {
        super(element, element, RESOLVER);
    }

    @Override
    protected PackageReference self() {
        return this;
    }

    @Override
    public PackageSolver newSolver() {
        return new PackageSolver(self());
    }

    @Override
    public void walkSolver(PackageSolver solver) {
        GoPsiScopesUtil.treeWalkUp(
                solver,
                getElement().getParent().getParent(),
                getElement().getContainingFile(),
                ResolveStates.initial());
    }

    @NotNull
    @Override
    public String getCanonicalText() {
        return getElement().getText();
    }

    @Override
    public boolean isReferenceTo(PsiElement element) {
        return false;
    }

    @NotNull
    @Override
    public Object[] getVariants() {
        return new Object[0];
    }

    @Override
    public boolean isSoft() {
        return false;
    }


}
