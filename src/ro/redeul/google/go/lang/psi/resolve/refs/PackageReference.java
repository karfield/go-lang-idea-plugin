package ro.redeul.google.go.lang.psi.resolve.refs;

import com.intellij.patterns.ElementPattern;
import com.intellij.psi.PsiElement;
import org.jetbrains.annotations.NotNull;
import ro.redeul.google.go.lang.psi.expressions.literals.GoLiteralIdentifier;
import ro.redeul.google.go.lang.psi.expressions.primary.GoLiteralExpression;
import ro.redeul.google.go.lang.psi.processors.ResolveStates;
import ro.redeul.google.go.lang.psi.resolve.GoReferenceWithSolver;
import ro.redeul.google.go.lang.psi.resolve.PackageSolver;
import ro.redeul.google.go.lang.psi.resolve.RefSolver;
import ro.redeul.google.go.lang.psi.resolve.ResolvingCache;
import ro.redeul.google.go.lang.psi.utils.GoPsiScopesUtil;

import static com.intellij.patterns.PlatformPatterns.psiElement;

public class PackageReference extends GoReferenceWithSolver<GoLiteralIdentifier, PackageSolver, PackageReference> {

    public static ElementPattern<GoLiteralIdentifier> MATCHER =
            psiElement(GoLiteralIdentifier.class)
                    .withParent(
                            psiElement(GoLiteralExpression.class)
                    );

    private static final ResolvingCache.Resolver<PackageReference> RESOLVER = ResolvingCache.makeDefault();

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
    public void walkSolver(RefSolver<?, ?> solver) {
        GoPsiScopesUtil.treeWalkUp(
                solver,
                getElement().getParent().getParent(),
                getElement().getContainingFile(),
                ResolveStates.initial());
    }


    @Override
    public boolean isReferenceTo(PsiElement element) {
        return false;
    }

}
