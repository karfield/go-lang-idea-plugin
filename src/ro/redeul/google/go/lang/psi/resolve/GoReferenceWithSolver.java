package ro.redeul.google.go.lang.psi.resolve;

import com.intellij.psi.PsiElement;
import org.jetbrains.annotations.NotNull;
import org.jetbrains.annotations.Nullable;
import ro.redeul.google.go.lang.psi.GoPsiElement;

public abstract class GoReferenceWithSolver<
        E extends GoPsiElement,
        S extends GoReferenceSolver<R, S>,
        R extends GoReferenceWithSolver<E, S, R>
        > extends GoReference<E, R> {

    protected abstract S newSolver(boolean collectAllVariants);

    protected abstract void walkSolver(S solver);

    @NotNull
    @Override
    public Object[] getVariants() {
        S solver = newSolver(true);
        walkSolver(solver);
        return solver.getVariants();
    }

    @Nullable
    @Override
    public PsiElement resolve() {
        return newSolver(false).resolve(self());
//        ResolveCache.getInstance(getElement().getProject()).resolveWithCaching(self(), resolver, true, false);
//        return null;
    }
}
